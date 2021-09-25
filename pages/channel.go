package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"codeberg.org/imabritishcow/librarian/api"
	"codeberg.org/imabritishcow/librarian/templates"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("Cache-Control", "public,max-age=1800")

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
