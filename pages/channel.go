package pages

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/imabritishcow/librarian/api"
	"github.com/imabritishcow/librarian/templates"
	"github.com/spf13/viper"
)


func ChannelHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
  w.WriteHeader(http.StatusOK)

	page := 1
	pageParam, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil {
		page = pageParam
	}

	claimType := "stream"
	if r.URL.Query().Get("claimType") != "" {
		claimType = r.URL.Query().Get("claimType")
	}

	orderBy := []string{"release_time"}
	if r.URL.Query().Get("orderBy") != "" {
		if r.URL.Query().Get("orderBy") == "trending" {
			orderBy = []string{"trending_group", "trending_mixed"}
		} else {
			orderBy = []string{r.URL.Query().Get("orderBy")}
		}
	}

	channelData := api.GetChannel(vars["channel"])
	followers, err := api.GetChannelFollowers(channelData.Id)
	if err != nil {
		log.Fatal(err)
	}
	videos := api.GetChannelVideos(page, channelData.Id, []string{claimType}, orderBy)

	channelTemplate, _ := template.ParseFS(templates.GetFiles(), "channel.html")
	err = channelTemplate.Execute(w, map[string]interface{}{
		"channel": channelData,
		"config": viper.AllSettings(),
		"followers": followers,
		"videos": videos,
		"query": map[string]interface{}{
			"page": page,
			"claimType": claimType,
			"orderBy": orderBy,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}