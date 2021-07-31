package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"codeberg.org/imabritishcow/librarian/types"
	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var waitingVideos sync.WaitGroup

func GetChannel(channel string) types.Channel {
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
	channelRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=resolve", "application/json", bytes.NewBuffer(resolveData))
	if err != nil {
		log.Fatal(err)
	}

	channelBody, err2 := ioutil.ReadAll(channelRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	channelData := gjson.Get(string(channelBody), "result."+channel)

	description := utils.ProcessText(channelData.Get("value.description").String(), true)
	thumbnail := channelData.Get("value.thumbnail.url").String()
	if thumbnail != "" {
		thumbnail = "/image?url="+thumbnail+"&hash="+utils.EncodeHMAC(thumbnail)
	}
	coverImg := channelData.Get("value.cover.url").String()
	if coverImg != "" {
		coverImg = "/image?url="+coverImg+"&hash="+utils.EncodeHMAC(coverImg)
	}

	return types.Channel{
		Name:        channelData.Get("name").String(),
		Title:       channelData.Get("value.title").String(),
		Id:          channelData.Get("claim_id").String(),
		Url:         strings.Replace(channelData.Get("canonical_url").String(), "lbry://", "https://"+viper.GetString("DOMAIN")+"/", 1),
		OdyseeUrl:   strings.ReplaceAll(channelData.Get("canonical_url").String(), "lbry://", "https://odysee.com/"),
		CoverImg:    coverImg,
		Description: template.HTML(description),
		Thumbnail:   thumbnail,
	}
}

func GetChannelFollowers(claimId string) (int64, error) {
	res, err := http.Get("https://api.lbry.com/subscription/sub_count?auth_token=" + viper.GetString("AUTH_TOKEN") + "&claim_id=" + claimId)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(res.Body)

	return gjson.Get(string(body), "data.0").Int(), err
}

func GetChannelVideos(page int, channelId string) []types.Video {
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
	channelDataRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=claim_search", "application/json", bytes.NewBuffer(channelData))
	if err != nil {
		log.Fatal(err)
	}

	channelDataBody, err := ioutil.ReadAll(channelDataRes.Body)
	if err != nil {
		log.Fatal(err)
	}

	videos := make([]types.Video, 0)
	videosData := gjson.Parse(string(channelDataBody))

	waitingVideos.Add(int(videosData.Get("result.items.#").Int()))
	videosData.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			go func() {
				claimId := value.Get("claim_id").String()
				lbryUrl := value.Get("canonical_url").String()
				channelLbryUrl := value.Get("signing_channel.canonical_url").String()

				time := time.Unix(value.Get("value.release_time").Int(), 0)
				thumbnail := value.Get("value.thumbnail.url").String()

				videos = append(videos, types.Video{
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
					ThumbnailUrl: "/image?url="+thumbnail+"&hash="+utils.EncodeHMAC(thumbnail),
					Views:        GetVideoViews(claimId),
					Timestamp:    time.Unix(),
					Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
					Duration:     utils.FormatDuration(value.Get("value.video.duration").Int()),
					RelTime:      humanize.Time(time),
				})
				waitingVideos.Done()
			}()

			return true
		},
	)
	waitingVideos.Wait()

	return videos
}
