package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"codeberg.org/imabritishcow/librarian/types"
	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var videoCache = cache.New(30*time.Minute, 15*time.Minute)

func GetVideo(channel string, video string, claimId string) (types.Video, error) {
	urls := []string{"lbry://" + channel + "/" + video}
	if channel == "" {
		urls = []string{"lbry://" + video + "#" + claimId}
	}

	cacheData, found := videoCache.Get(urls[0])
	if found {
		return cacheData.(types.Video), nil
	}

	resolveDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "resolve",
		"params": map[string]interface{}{
			"urls":                     urls,
			"include_purchase_receipt": true,
			"include_is_my_output":     true,
		},
		"id": time.Now().Unix(),
	}
	resolveData, _ := json.Marshal(resolveDataMap)
	videoDataRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=resolve", "application/json", bytes.NewBuffer(resolveData))

	videoDataBody, err := ioutil.ReadAll(videoDataRes.Body)

	videoData := gjson.Get(string(videoDataBody), "result.lbry*")
	if videoData.Get("error.name").String() != "" {
		return types.Video{}, fmt.Errorf("API Error: " + videoData.Get("error.name").String() + videoData.Get("error.text").String())
	}

	returnData := ProcessVideo(videoData)
	videoCache.Set(urls[0], returnData, cache.DefaultExpiration)
	return returnData, err
}

func GetVideoViews(claimId string) int64 {
	cacheData, found := videoCache.Get(claimId + "-views")
	if found {
		return cacheData.(int64)
	}

	viewCountRes, err := http.Get("https://api.odysee.com/file/view_count?auth_token=" + viper.GetString("AUTH_TOKEN") + "&claim_id=" + claimId)
	if err != nil {
		fmt.Println(err)
	}

	viewCountBody, err2 := ioutil.ReadAll(viewCountRes.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	returnData := gjson.Get(string(viewCountBody), "data.0").Int()
	videoCache.Set(claimId + "-views", returnData, cache.DefaultExpiration)
	return returnData
}

func GetLikeDislike(claimId string) []int64 {
	cacheData, found := videoCache.Get(claimId + "-reactions")
	if found {
		return cacheData.([]int64)
	}

	likeDislikeRes, err := http.PostForm("https://api.odysee.com/reaction/list", url.Values{
		"claim_ids": []string{claimId},
	})
	if err != nil {
		fmt.Println(err)
	}

	likeDislikeBody, err2 := ioutil.ReadAll(likeDislikeRes.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	returnData := []int64{
		gjson.Get(string(likeDislikeBody), "data.others_reactions."+claimId+".like").Int(),
		gjson.Get(string(likeDislikeBody), "data.others_reactions."+claimId+".dislike").Int(),
	}
	videoCache.Set(claimId + "-reactions", returnData, cache.DefaultExpiration)
	return returnData
}

func GetVideoStream(video string) string {
	cacheData, found := videoCache.Get(video + "-strean")
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
	videoStreamRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=get", "application/json", bytes.NewBuffer(getData))
	if err != nil {
		fmt.Println(err)
	}

	videoStreamBody, err2 := ioutil.ReadAll(videoStreamRes.Body)
	if err2 != nil {
		fmt.Println(err2)
	}
	
	returnData := gjson.Get(string(videoStreamBody), "result.streaming_url").String()
	if viper.GetString("CDN_VIDEO_PROXY") != "" {
		returnData = strings.ReplaceAll(returnData, "cdn.lbryplayer.xyz", viper.GetString("CDN_VIDEO_PROXY"))
	}
	
	videoCache.Set(video + "-stream", returnData, cache.DefaultExpiration)
	return returnData
}

func ProcessVideo(videoData gjson.Result) types.Video {
	tags := make([]string, 0)
	videoData.Get("value.tags").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			tags = append(tags, value.String())
			return true
		},
	)

	claimId := videoData.Get("claim_id").String()
	lbryUrl := videoData.Get("canonical_url").String()
	channelLbryUrl := videoData.Get("signing_channel.canonical_url").String()

	time := time.Unix(videoData.Get("value.release_time").Int(), 0)
	thumbnail := videoData.Get("value.thumbnail.url").String()
	channelThumbnail := videoData.Get("signing_channel.value.thumbnail.url").String()
	if channelThumbnail != "" {
		channelThumbnail = "/image?url=" + channelThumbnail + "&hash=" + utils.EncodeHMAC(channelThumbnail)
	}

	likeDislike := GetLikeDislike(claimId)
	
	return types.Video{
		Url:       utils.LbryTo(lbryUrl, "http"),
		LbryUrl:   lbryUrl,
		RelUrl:    utils.LbryTo(lbryUrl, "rel"),
		OdyseeUrl: utils.LbryTo(lbryUrl, "odysee"),
		ClaimId:   videoData.Get("claim_id").String(),
		Channel: types.Channel{
			Name:        videoData.Get("signing_channel.name").String(),
			Title:       videoData.Get("signing_channel.value.title").String(),
			Id:          videoData.Get("signing_channel.claim_id").String(),
			Url:         utils.LbryTo(channelLbryUrl, "http"),
			RelUrl:      utils.LbryTo(channelLbryUrl, "rel"),
			OdyseeUrl:   utils.LbryTo(channelLbryUrl, "odysee"),
			Description: template.HTML(videoData.Get("signing_channel.value.description").String()),
			Thumbnail:   channelThumbnail,
		},
		Duration:     utils.FormatDuration(videoData.Get("value.video.duration").Int()),
		Title:        videoData.Get("value.title").String(),
		ThumbnailUrl: "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
		Description:  template.HTML(utils.ProcessText(videoData.Get("value.description").String(), true)),
		License:      videoData.Get("value.license").String(),
		Views:        GetVideoViews(claimId),
		Likes:        likeDislike[0],
		Dislikes:     likeDislike[1],
		Tags:         tags,
		RelTime:      humanize.Time(time),
		Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
	}
}
