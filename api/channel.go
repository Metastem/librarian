package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var channelCache = cache.New(30*time.Minute, 15*time.Minute)

func GetChannel(channel string, getFollowers bool) (types.Channel, error) {
	cacheData, found := channelCache.Get(channel)
	if found {
		return cacheData.(types.Channel), nil
	}

	Client := utils.NewClient()
	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     []string{channel},
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)
	channelRes, err := Client.Post(viper.GetString("API_URL")+"?m=resolve", "application/json", bytes.NewBuffer(resolveData))
	if err != nil {
		return types.Channel{}, err
	}

	channelBody, err := ioutil.ReadAll(channelRes.Body)
	if err != nil {
		return types.Channel{}, err
	}

	data := gjson.Get(string(channelBody), "result."+strings.ReplaceAll(channel, ".", "\\."))
	channelData, err := ProcessChannel(data, getFollowers)
	if err != nil {
		return types.Channel{}, err
	}
	
	channelCache.Set(channel, channelData, cache.DefaultExpiration)
	return channelData, nil
}

func ProcessChannel(data gjson.Result, getFollowers bool) (types.Channel, error) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	description := ""
	thumbnail := data.Get("value.thumbnail.url").String()
	go func() {
		defer wg.Done()
		description = utils.ProcessText(data.Get("value.description").String(), true)
		if thumbnail != "" {
			thumbnail = "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail)
		}
	}()

	coverImg := data.Get("value.cover.url").String()
	go func() {
		defer wg.Done()
		if coverImg != "" {
			coverImg = "/image?url=" + coverImg + "&hash=" + utils.EncodeHMAC(coverImg)
		}
	}()

	followers, err := int64(0), error(nil)
	go func() {
		defer wg.Done()
		if getFollowers {
			followers, err = GetChannelFollowers(data.Get("claim_id").String())
		}
	}()
	if err != nil {
		return types.Channel{}, err
	}

	wg.Wait()

	url, err := utils.LbryTo(data.Get("canonical_url").String())
	if err != nil {
		return types.Channel{}, err
	}

	return types.Channel{
		Name:           data.Get("name").String(),
		Title:          data.Get("value.title").String(),
		Id:             data.Get("claim_id").String(),
		Url:            url["http"],
		OdyseeUrl:      url["odysee"],
		RelUrl:         url["rel"],
		CoverImg:       coverImg,
		Description:    template.HTML(description),
		DescriptionTxt: bluemonday.StrictPolicy().Sanitize(description),
		Thumbnail:      thumbnail,
		Followers:      followers,
		ValueType:      data.Get("value_type").String(),
		UploadCount:    data.Get("meta.claims_in_channel").Int(),
	}, nil
}

func GetChannelFollowers(claimId string) (int64, error) {
	cacheData, found := channelCache.Get(claimId + "-followers")
	if found {
		return cacheData.(int64), nil
	}

	Client := utils.NewClient()
	res, err := Client.Get("https://api.odysee.com/subscription/sub_count?auth_token=" + viper.GetString("AUTH_TOKEN") + "&claim_id=" + claimId)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(res.Body)

	returnData := gjson.Get(string(body), "data.0").Int()
	channelCache.Set(claimId+"-followers", returnData, cache.DefaultExpiration)
	return returnData, err
}

func GetChannelClaims(page int, channelId string) ([]types.Claim, error) {
	cacheData, found := channelCache.Get(channelId + "-claims-" + fmt.Sprint(page))
	if found {
		return cacheData.([]types.Claim), nil
	}

	Client := utils.NewClient()
	channelDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":                20,
			"page":                     page,
			"no_totals":                true,
			"claim_type":               []string{"stream"},
			"order_by":                 []string{"release_time"},
			"fee_amount":               "<=0",
			"channel_ids":              []string{channelId},
			"release_time":             "<" + fmt.Sprint(time.Now().Unix()),
			"include_purchase_receipt": true,
		},
	}
	channelData, _ := json.Marshal(channelDataMap)
	channelDataRes, err := Client.Post(viper.GetString("API_URL")+"?m=claim_search", "application/json", bytes.NewBuffer(channelData))
	if err != nil {
		return []types.Claim{}, nil
	}

	channelDataBody, err := ioutil.ReadAll(channelDataRes.Body)
	if err != nil {
		return []types.Claim{}, nil
	}

	claims := make([]types.Claim, 0)
	claimsData := gjson.Parse(string(channelDataBody))

	wg := sync.WaitGroup{}
	claimsData.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)

			go func() {
				defer wg.Done()

				claimId := value.Get("claim_id").String()
				lbryUrl := value.Get("canonical_url").String()
				channelLbryUrl := value.Get("signing_channel.canonical_url").String()

				time := time.Unix(value.Get("value.release_time").Int(), 0)
				thumbnail := value.Get("value.thumbnail.url").String()
				thumbnail = url.QueryEscape(thumbnail)

				views, err := GetViews(claimId)
				if err != nil {
					fmt.Println(err)
					return
				}

				url, err := utils.LbryTo(lbryUrl)
				if err != nil {
					fmt.Println(err)
					return
				}

				channelUrl, err := utils.LbryTo(channelLbryUrl)
				if err != nil {
					fmt.Println(err)
					return
				}

				claims = append(claims, types.Claim{
					Url:       url["http"],
					LbryUrl:   lbryUrl,
					RelUrl:    url["rel"],
					OdyseeUrl: url["odysee"],
					ClaimId:   value.Get("claim_id").String(),
					Channel: types.Channel{
						Name:      value.Get("signing_channel.name").String(),
						Title:     value.Get("signing_channel.value.title").String(),
						Id:        value.Get("signing_channel.claim_id").String(),
						Url:       channelUrl["http"],
						RelUrl:    channelUrl["rel"],
						OdyseeUrl: channelUrl["odysee"],
					},
					Description:  template.HTML(utils.ProcessText(value.Get("value.description").String(), true)),
					Title:        value.Get("value.title").String(),
					ThumbnailUrl: "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
					Views:        views,
					Timestamp:    time.Unix(),
					Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
					Duration:     utils.FormatDuration(value.Get("value.video.duration").Int()),
					RelTime:      humanize.Time(time),
					MediaType:    value.Get("value.source.media_type").String(),
					StreamType:   value.Get("value.stream_type").String(),
					SrcSize:      value.Get("value.source.size").String(),
				})
			}()

			return true
		},
	)
	wg.Wait()

	channelCache.Set(channelId+"-claims-"+fmt.Sprint(page), claims, cache.DefaultExpiration)
	return claims, nil
}
