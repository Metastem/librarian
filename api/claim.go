package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var claimCache = cache.New(30*time.Minute, 15*time.Minute)

func GetClaim(channel string, video string, claimId string) (types.Claim, error) {
	urls := []string{"lbry://" + channel + "/" + video}
	if channel == "" {
		urls = []string{"lbry://" + video + "#" + claimId}
	}

	cacheData, found := claimCache.Get(urls[0])
	if found {
		return cacheData.(types.Claim), nil
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
	claimDataRes, err := http.Post(viper.GetString("API_URL")+"?m=resolve", "application/json", bytes.NewBuffer(resolveData))

	claimDataBody, err := ioutil.ReadAll(claimDataRes.Body)

	claimData := gjson.Get(string(claimDataBody), "result.lbry*")
	if claimData.Get("error.name").String() != "" {
		return types.Claim{}, fmt.Errorf("API Error: " + claimData.Get("error.name").String() + claimData.Get("error.text").String())
	}

	returnData := ProcessClaim(claimData)
	claimCache.Set(urls[0], returnData, cache.DefaultExpiration)
	return returnData, err
}

func ProcessClaim(claimData gjson.Result) types.Claim {
	tags := make([]string, 0)
	claimData.Get("value.tags").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			tags = append(tags, value.String())
			return true
		},
	)

	claimId := claimData.Get("claim_id").String()
	lbryUrl := claimData.Get("canonical_url").String()
	channelLbryUrl := claimData.Get("signing_channel.canonical_url").String()

	time := time.Unix(claimData.Get("value.release_time").Int(), 0)
	thumbnail := claimData.Get("value.thumbnail.url").String()
	channelThumbnail := claimData.Get("signing_channel.value.thumbnail.url").String()
	if channelThumbnail != "" {
		channelThumbnail = "/image?url=" + channelThumbnail + "&hash=" + utils.EncodeHMAC(channelThumbnail)
	}

	likeDislike := GetLikeDislike(claimId)

	return types.Claim{
		Url:       utils.LbryTo(lbryUrl, "http"),
		LbryUrl:   lbryUrl,
		RelUrl:    utils.LbryTo(lbryUrl, "rel"),
		OdyseeUrl: utils.LbryTo(lbryUrl, "odysee"),
		ClaimId:   claimData.Get("claim_id").String(),
		Channel: types.Channel{
			Name:        claimData.Get("signing_channel.name").String(),
			Title:       claimData.Get("signing_channel.value.title").String(),
			Id:          claimData.Get("signing_channel.claim_id").String(),
			Url:         utils.LbryTo(channelLbryUrl, "http"),
			RelUrl:      utils.LbryTo(channelLbryUrl, "rel"),
			OdyseeUrl:   utils.LbryTo(channelLbryUrl, "odysee"),
			Description: template.HTML(claimData.Get("signing_channel.value.description").String()),
			Thumbnail:   channelThumbnail,
		},
		Duration:     utils.FormatDuration(claimData.Get("value.video.duration").Int()),
		Title:        claimData.Get("value.title").String(),
		ThumbnailUrl: "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail),
		Description:  template.HTML(utils.ProcessText(claimData.Get("value.description").String(), true)),
		License:      claimData.Get("value.license").String(),
		Views:        GetViews(claimId),
		Likes:        likeDislike[0],
		Dislikes:     likeDislike[1],
		Tags:         tags,
		RelTime:      humanize.Time(time),
		Date:         time.Month().String() + " " + fmt.Sprint(time.Day()) + ", " + fmt.Sprint(time.Year()),
		StreamType: 	claimData.Get("value.stream_type").String(),
	}
}

func GetViews(claimId string) int64 {
	cacheData, found := claimCache.Get(claimId + "-views")
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
	claimCache.Set(claimId+"-views", returnData, cache.DefaultExpiration)
	return returnData
}

func GetLikeDislike(claimId string) []int64 {
	cacheData, found := claimCache.Get(claimId + "-reactions")
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
	claimCache.Set(claimId+"-reactions", returnData, cache.DefaultExpiration)
	return returnData
}