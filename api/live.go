package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func GetLive(claimId string) (types.Live, error) {
	liveRes, err := http.Get("https://api.odysee.live/livestream/is_live?channel_claim_id=" + claimId)
	if err != nil {
		return types.Live{}, err
	}

	liveBody, err := ioutil.ReadAll(liveRes.Body)
	if err != nil {
		return types.Live{}, err
	}

	data := gjson.Parse(string(liveBody))
	if !data.Get("success").Bool() {
		return types.Live{}, fmt.Errorf(data.Get("error").String())
	}

	timestamp, err := time.Parse("2006-01-02T15:04:05.999Z", data.Get("data.Start").String())
	if err != nil {
		return types.Live{}, err
	}

	thumbnail := data.Get("data.ThumbnailURL").String()
	thumbnail = url.QueryEscape(thumbnail)
	thumbnail = "/image?url=" + thumbnail + "&hash=" + utils.EncodeHMAC(thumbnail)

	streamUrl := strings.ReplaceAll(data.Get("data.VideoURL").String(), "https://cloud.odysee.live", "/live")
	if viper.GetString("LIVE_STREAMING_URL") != "" {
		streamUrl = strings.ReplaceAll(data.Get("data.VideoURL").String(), "https://cloud.odysee.live", viper.GetString("LIVE_STREAMING_URL"))	
	}

	return types.Live{
		RelTime: humanize.Time(timestamp),
		Time: timestamp.Format("Jan 2, 2006 03:04 PM"),
		ThumbnailUrl: thumbnail,
		StreamUrl: streamUrl,
		Live: data.Get("data.Live").Bool(),
	}, nil
}