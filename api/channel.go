package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

type Channel struct {
	Name string
	Title string
	Id string
	Url string
	OdyseeUrl string
	CoverImg string
	Description string
	Thumbnail string
}

func GetChannel(channel string) Channel {
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

	return Channel{
		Name: channelData.Get("name").String(),
		Title: channelData.Get("value.title").String(),
		Id: channelData.Get("claim_id").String(),
		Url: strings.Replace(channelData.Get("canonical_url").String(), "lbry://", "https://"+viper.GetString("DOMAIN")+"/", 1),
		OdyseeUrl: strings.ReplaceAll(channelData.Get("canonical_url").String(), "lbry://", "https://odysee.com/"),
		CoverImg: channelData.Get("value.cover.url").String(),
		Description: channelData.Get("value.description").String(),
		Thumbnail: channelData.Get("value.thumbnail.url").String(),
	}
}
