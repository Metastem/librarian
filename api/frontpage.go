package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/data"
	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var fpCache = cache.New(30*time.Minute, 30*time.Minute)

func GetFrontpageVideos() ([]types.Claim, error) {
	cacheData, found := fpCache.Get("fp")
	if found {
		return cacheData.([]types.Claim), nil
	}

	Client := utils.NewClient()
	claimSearchData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":                20,
			"no_totals":                true,
			"claim_type":               "stream",
			"any_tags":                 []string{},
			"not_tags":                 []string{"porn", "porno", "nsfw", "mature", "xxx", "sex", "creampie", "blowjob", "handjob", "vagina", "boobs", "big boobs", "big dick", "pussy", "cumshot", "anal", "hard fucking", "ass", "fuck", "hentai"},
			"channel_ids":              data.Home,
			"not_channel_ids":          []string{},
			"order_by":                 []string{"release_time"},
			"fee_amount":               "<=0",
			"release_time":             ">" + fmt.Sprint(time.Now().Unix()-15778458),
			"include_purchase_receipt": true,
		},
	}
	claimSearchReqData, _ := json.Marshal(claimSearchData)
	frontpageDataRes, err := Client.Post(viper.GetString("API_URL")+"?m=claim_search", "application/json", bytes.NewBuffer(claimSearchReqData))
	if err != nil {
		return []types.Claim{}, err
	}

	frontpageDataBody, err := ioutil.ReadAll(frontpageDataRes.Body)
	if err != nil {
		return []types.Claim{}, err
	}

	claims := make([]types.Claim, 0)
	claimsData := gjson.Parse(string(frontpageDataBody))

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

				claims = append(claims, types.Claim{
					Url:       utils.LbryTo(lbryUrl, "http"),
					LbryUrl:   lbryUrl,
					RelUrl:    utils.LbryTo(lbryUrl, "rel"),
					OdyseeUrl: utils.LbryTo(lbryUrl, "odysee"),
					ClaimId:   value.Get("claim_id").String(),
					Channel: types.Channel{
						Name:      value.Get("signing_channel.name").String(),
						Title:     value.Get("signing_channel.value.title").String(),
						Id:        value.Get("signing_channel.claim_id").String(),
						Url:       utils.LbryTo(channelLbryUrl, "http"),
						RelUrl:    utils.LbryTo(channelLbryUrl, "rel"),
						OdyseeUrl: utils.LbryTo(channelLbryUrl, "odysee"),
					},
					Title:        value.Get("value.title").String(),
					ThumbnailUrl: "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
					Views:        GetViews(claimId),
					Timestamp:    time.Unix(),
					Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
					Duration:     utils.FormatDuration(value.Get("value.video.duration").Int()),
					RelTime:      humanize.Time(time),
					StreamType:   value.Get("value.stream_type").String(),
				})
			}()

			return true
		},
	)
	wg.Wait()

	commentCache.Set("fp", claims, cache.DefaultExpiration)
	return claims, nil
}
