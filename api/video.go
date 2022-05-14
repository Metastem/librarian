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

func GetVideoStream(video string) (string, error) {
	cacheData, found := videoCache.Get(video + "-stream")
	if found {
		return cacheData.(string), nil
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
	videoStreamRes, err := http.Post(viper.GetString("STREAMING_API_URL")+"?m=get", "application/json", bytes.NewBuffer(getData))
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
		returnData = strings.ReplaceAll(returnData, "https://player.odycdn.com", viper.GetString("VIDEO_STREAMING_URL"))
	}

	videoCache.Set(video+"-stream", returnData, cache.DefaultExpiration)
	return returnData, nil
}

func CheckHLS(url string) (string, error) {
	res, err := http.Head(url)
	if err != nil {
		return "", err
	}
	if res.StatusCode == 403 {
		return "", fmt.Errorf("this content cannot be accessed due to a DMCA request")
	}
	if res.Header.Get("Content-Type") == "application/x-mpegurl" {
		return res.Request.URL.String(), nil
	}
	return res.Header.Get("Content-Type"), nil
}