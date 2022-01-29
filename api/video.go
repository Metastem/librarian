package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"codeberg.org/librarian/librarian/utils"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var videoCache = cache.New(30*time.Minute, 15*time.Minute)

func GetVideoStream(video string) (string, error) {
	cacheData, found := videoCache.Get(video + "-stream")
	if found {
		return cacheData.(string), nil
	}

	Client := utils.NewClient()
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
	videoStreamRes, err := Client.Post(viper.GetString("STREAMING_API_URL")+"?m=get", "application/json", bytes.NewBuffer(getData))
	if err != nil {
		return "", err
	}

	videoStreamBody, err := ioutil.ReadAll(videoStreamRes.Body)
	if err != nil {
		return "", err
	}

	returnData := gjson.Get(string(videoStreamBody), "result.streaming_url").String()
	if viper.GetString("VIDEO_STREAMING_URL") != "" {
		returnData = strings.ReplaceAll(returnData, "http://localhost:5280", viper.GetString("VIDEO_STREAMING_URL"))
		returnData = strings.ReplaceAll(returnData, "https://cdn.lbryplayer.xyz", viper.GetString("VIDEO_STREAMING_URL"))
	}

	videoCache.Set(video+"-stream", returnData, cache.DefaultExpiration)
	return returnData, nil
}

func GetVideoStreamType(url string) (string, error) {
	res, err := http.Head(url)
	if err != nil {
		return "", err
	}
	return res.Header.Get("Content-Type"), nil
}