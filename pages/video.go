package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"codeberg.org/imabritishcow/librarian/api"
	"codeberg.org/imabritishcow/librarian/templates"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Cache-Control", "public,max-age=3600")

	videoData := api.GetVideo(vars["channel"], vars["video"], "")

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), videoData.ClaimId) {
		blockTemplate, _ := template.ParseFS(templates.GetFiles(), "blocked.html")
		err := blockTemplate.Execute(w, map[string]interface{}{
			"video": videoData,
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	videoStream := api.GetVideoStream(videoData.LbryUrl)
	comments := api.GetComments(videoData.ClaimId, videoData.Channel.Id, videoData.Channel.Name)

	videoTemplate, _ := template.ParseFS(templates.GetFiles(), "video.html")
	err := videoTemplate.Execute(w, map[string]interface{}{
		"stream":         videoStream,
		"video":          videoData,
		"comments":       comments,
		"commentsLength": len(comments),
		"config":         viper.AllSettings(),
	})
	if err != nil {
		fmt.Println(err)
	}
}
