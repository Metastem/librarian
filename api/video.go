package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var videoCache = cache.New(30*time.Minute, 15*time.Minute)

func GetVideoStream(video string) string {
	cacheData, found := videoCache.Get(video + "-stream")
	if found {
		return cacheData.(string)
	}

	getDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "get",
		"params": map[string]interface{}{
			"uri":       video,
			"save_file": false,
		},
		"id": time.Now().Unix(),
	}
	getData, _ := json.Marshal(getDataMap)
	videoStreamRes, err := Client.Post(viper.GetString("API_URL")+"?m=get", "application/json", bytes.NewBuffer(getData))
	if err != nil {
		fmt.Println(err)
	}

	videoStreamBody, err2 := ioutil.ReadAll(videoStreamRes.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	returnData := gjson.Get(string(videoStreamBody), "result.streaming_url").String()
	if viper.GetString("VIDEO_STREAMING_URL") != "" {
		returnData = strings.ReplaceAll(returnData, "http://localhost:5280", viper.GetString("VIDEO_STREAMING_URL"))
		returnData = strings.ReplaceAll(returnData, "https://cdn.lbryplayer.xyz", viper.GetString("VIDEO_STREAMING_URL"))
	}

	videoCache.Set(video+"-stream", returnData, cache.DefaultExpiration)
	return returnData
}

func GetVideoStreamType(url string) string {
	res, _ := http.Head(url)
	return res.Header.Get("Content-Type")
}

func GetStcStream(claimId string) map[string]string {
	stcRes, err := http.Get(viper.GetString("STC_URL") + "/find?claim_id=" + claimId)
	if err != nil {
		fmt.Println(err)
	}

	stcBody, err2 := ioutil.ReadAll(stcRes.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	data := gjson.Parse(string(stcBody))

	return map[string]string{
		"fhd": data.Get("1080p").String(),
		"hd": data.Get("720p").String(),
		"sd": data.Get("480p").String(),
	}
}