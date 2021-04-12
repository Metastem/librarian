package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
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
	videoDataReq, err := http.NewRequest(http.MethodPost, "https://api.lbry.tv/api/v1/proxy?m=resolve", bytes.NewBuffer(resolveData))
	if err != nil {
		log.Fatal(err)
	}
	videoDataReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0")
	videoDataReq.Header.Set("Content-Type", "application/json")

	getDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "get",
		"params": map[string]interface{}{
			"uri":      "lbry://" + channel + "/" + video,
			"save_file": false,
		},
		"id": time.Now().Unix(),
	}
	getData, _ := json.Marshal(getDataMap)
	videoStreamReq, err := http.NewRequest(http.MethodPost, "https://api.lbry.tv/api/v1/proxy?m=get", bytes.NewBuffer(getData))
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
	videoStream := gjson.Get(string(videoStreamBody), "result.streaming_url")
	
	videos := make([]map[string]interface{}, 0)
	videoData.ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			videos = append(videos, map[string]interface{}{
				"url":          value.Get("permanent_url").String(),
				"channel":      value.Get("signing_channel.name").String(),
				"channelPfp":   value.Get("signing_channel.value.cover.url").String(),
				"title":        value.Get("value.title").String(),
				"thumbnailUrl": value.Get("value.thumbnail.url").String(),
				"description":  value.Get("value.description").String(),
				"license":      value.Get("value.license").String(),
			})

			return true
		},
	)

	return VideoResult{
		StreamUrl: videoStream.String(),
		Videos:    videos,
	}
}
