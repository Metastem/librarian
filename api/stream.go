package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var streamCache = cache.New(30*time.Minute, 15*time.Minute)

func GetStream(video string) (types.Stream, error) {
	cacheData, found := streamCache.Get(video + "-stream")
	if found {
		return cacheData.(types.Stream), nil
	}

	reqDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "get",
		"params": map[string]interface{}{
			"uri":       video,
			"save_file": false,
		},
		"id": time.Now().Unix(),
	}
	reqData, err := json.Marshal(reqDataMap)
	if err != nil {
		return types.Stream{}, err
	}

	data, err := utils.RequestJSON(viper.GetString("STREAMING_API_URL")+"?m=get", bytes.NewBuffer(reqData), true)
	if err != nil {
		return types.Stream{}, err
	}

	streamUrl := data.Get("result.streaming_url").String()
	streamUrl = strings.ReplaceAll(streamUrl, "source.odycdn.com", "player.odycdn.com")
	if viper.GetString("VIDEO_STREAMING_URL") != "" {
		streamUrl = strings.ReplaceAll(streamUrl, "http://localhost:5280", viper.GetString("VIDEO_STREAMING_URL"))
		streamUrl = strings.ReplaceAll(streamUrl, "https://player.odycdn.com", viper.GetString("VIDEO_STREAMING_URL"))
	}

	stream, err := checkStream(streamUrl)
	if err != nil {
		return types.Stream{}, err
	}

	streamCache.Set(video+"-stream", stream, cache.DefaultExpiration)
	return stream, nil
}

func checkStream(url string) (types.Stream, error) {
	res, err := http.Head(url)
	if err != nil {
		return types.Stream{}, err
	}

	if res.StatusCode == 403 {
		return types.Stream{}, fmt.Errorf("this content cannot be accessed due to a DMCA request")
	}

	isHls := res.Header.Get("Content-Type") == "application/x-mpegurl"
	return types.Stream{
		Type: res.Header.Get("Content-Type"),
		URL: res.Request.URL.String(),
		HLS: isHls,
	}, nil
}