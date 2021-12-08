package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/templates"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("Cache-Control", "public,max-age=1800")
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("Referrer-Policy", "no-referrer")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("Strict-Transport-Security", "max-age=31557600")
	w.Header().Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=(), battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	w.Header().Add("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	page := 1
	pageParam, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channelData := api.GetChannel(vars["channel"], true)

	if (channelData.Id == "") {
		notFoundTemplate, _ := template.ParseFS(templates.GetFiles(), "404.html")
		err := notFoundTemplate.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	/*TO-DO: Add playlists

	videos := make([]types.Video, 0)
	claimType := r.URL.Query().Get("claimType")
	if claimType == "" || claimType == "stream" {
		claimType = "stream"
		videos := api.GetChannelVideos(page, channelData.Id)
		sort.Slice(videos, func (i int, j int) bool {
			return videos[i].Timestamp > videos[j].Timestamp
		})
	}*/

	videos := api.GetChannelVideos(page, channelData.Id)
	sort.Slice(videos, func(i int, j int) bool {
		return videos[i].Timestamp > videos[j].Timestamp
	})

	channelTemplate, _ := template.ParseFS(templates.GetFiles(), "channel.html")
	err = channelTemplate.Execute(w, map[string]interface{}{
		"channel":   channelData,
		"config":    viper.AllSettings(),
		"videos":    videos,
		"query": map[string]interface{}{
			"page":      fmt.Sprint(page),
			"nextPage":  fmt.Sprint(page + 1),
			"prevPage":  fmt.Sprint(page - 1),
			"page0":     "0",
			"claimType": "stream",
			"stream":    "stream",
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}
