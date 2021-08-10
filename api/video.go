package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"codeberg.org/imabritishcow/librarian/types"
	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func GetVideo(channel string, video string) types.Video {
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
	videoDataRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=resolve", "application/json", bytes.NewBuffer(resolveData))
	if err != nil {
		log.Fatal(err)
	}

	videoDataBody, err2 := ioutil.ReadAll(videoDataRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	videoData := gjson.Get(string(videoDataBody), "result.lbry*")

	return ProcessVideo(videoData)
}

func GetVideoViews(claimId string) int64 {
	viewCountRes, err := http.Get("https://api.odysee.com/file/view_count?auth_token=" + viper.GetString("AUTH_TOKEN") + "&claim_id=" + claimId)
	if err != nil {
		log.Fatal(err)
	}

	viewCountBody, err2 := ioutil.ReadAll(viewCountRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	println(string(viewCountBody))

	return gjson.Get(string(viewCountBody), "data.0").Int()
}

func GetLikeDislike(claimId string) []int64 {
	likeDislikeRes, err := http.PostForm("https://api.odysee.com/reaction/list", url.Values{
		"claim_ids": []string{claimId},
	})
	if err != nil {
		log.Fatal(err)
	}

	likeDislikeBody, err2 := ioutil.ReadAll(likeDislikeRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	return []int64{
		gjson.Get(string(likeDislikeBody), "data.others_reactions."+claimId+".like").Int(),
		gjson.Get(string(likeDislikeBody), "data.others_reactions."+claimId+".dislike").Int(),
	}
}

func GetVideoStream(video string) string {
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
		log.Fatal(err)
	}

	videoStreamBody, err2 := ioutil.ReadAll(videoStreamRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	return gjson.Get(string(videoStreamBody), "result.streaming_url").String()
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
		Title:        videoData.Get("value.title").String(),
		ThumbnailUrl: "/image?url" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
		Description:  template.HTML(utils.ProcessText(videoData.Get("value.description").String(), true)),
		License:      videoData.Get("value.license").String(),
		Views:        GetVideoViews(claimId),
		Likes:        likeDislike[0],
		Dislikes:     likeDislike[1],
		Tags:         tags,
		Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
	}
}
