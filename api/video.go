package api

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"mvdan.cc/xurls/v2"
)

type VideoResult struct {
	StreamUrl string
	Videos    []map[string]interface{}
}

func GetVideo(channel string, video string) VideoResult {
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}

	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     []string{"lbry://" + channel + "/" + video},
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)
	videoDataReq, err := http.NewRequest(http.MethodPost, viper.GetString("API_URL")+"/api/v1/proxy?m=resolve", bytes.NewBuffer(resolveData))
	if err != nil {
		log.Fatal(err)
	}
	videoDataReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0")
	videoDataReq.Header.Set("Content-Type", "application/json")

	getDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "get",
		"params": map[string]interface{}{
			"uri":       "lbry://" + channel + "/" + video,
			"save_file": false,
		},
		"id": time.Now().Unix(),
	}
	getData, _ := json.Marshal(getDataMap)
	videoStreamReq, err := http.NewRequest(http.MethodPost, viper.GetString("API_URL")+"/api/v1/proxy?m=get", bytes.NewBuffer(getData))
	if err != nil {
		log.Fatal(err)
	}
	videoStreamReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0")
	videoStreamReq.Header.Set("Content-Type", "application/json")

	videoDataRes, getErr1 := httpClient.Do(videoDataReq)
	if getErr1 != nil {
		log.Fatal(getErr1)
	}

	videoDataBody, readErr1 := ioutil.ReadAll(videoDataRes.Body)
	if err != nil {
		log.Fatal(readErr1)
	}

	videoStreamRes, getErr2 := httpClient.Do(videoStreamReq)
	if getErr2 != nil {
		log.Fatal(getErr2)
	}

	videoStreamBody, readErr2 := ioutil.ReadAll(videoStreamRes.Body)
	if err != nil {
		log.Fatal(readErr2)
	}

	videoData := gjson.Get(string(videoDataBody), "result")
	videoStream := gjson.Get(string(videoStreamBody), "result.streaming_url").String()

	videos := make([]map[string]interface{}, 0)
	videoData.ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			tags := make([]string, 0)
			value.Get("value.tags").ForEach(
				func(key gjson.Result, value gjson.Result) bool {
					tags = append(tags, value.String())
					return true
				},
			)

			p := bluemonday.UGCPolicy()

			description := value.Get("value.description").String()
			description = p.Sanitize(description)
			description = strings.ReplaceAll(description, "\n", "<br>")
			description = xurls.Relaxed().ReplaceAllString(description, "<a href=\"$1$3$4\">$1$3$4</a>")

			channelId := value.Get("signing_channel.short_url").String()
			channelId = strings.ReplaceAll(channelId, "lbry://", "")
			channelId = strings.ReplaceAll(channelId, "#", ":")

			videos = append(videos, map[string]interface{}{
				"url":          strings.Replace(value.Get("canonical_url").String(), "lbry://", "https://"+viper.GetString("DOMAIN")+"/", 1),
				"channel":      map[string]interface{}{
					"name": value.Get("signing_channel.value.title").String(),
					"id": channelId,
					"pfp": value.Get("signing_channel.value.cover.url").String(),
				},
				"tags":         tags,
				"title":        value.Get("value.title").String(),
				"thumbnailUrl": template.URL(value.Get("value.thumbnail.url").String()),
				"description":  template.HTML(description),
				"license":      value.Get("value.license").String(),
				"video": map[string]interface{}{
					"duration": value.Get("value.video.duration").Int(),
					"height":   value.Get("value.video.height").Int(),
					"width":    value.Get("value.video.width").Int(),
				},
			})

			return true
		},
	)

	return VideoResult{
		StreamUrl: videoStream,
		Videos:    videos,
	}
}
