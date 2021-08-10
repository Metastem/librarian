package pages

import (
	"html/template"
	"log"
	"net/http"
	"sort"

	"codeberg.org/imabritishcow/librarian/api"
	"codeberg.org/imabritishcow/librarian/templates"
	"github.com/spf13/viper"
)

func FrontpageHandler(w http.ResponseWriter, r *http.Request) {
	videos := api.GetFrontpageVideos()
	sort.Slice(videos, func(i int, j int) bool {
		return videos[i].Timestamp > videos[j].Timestamp
	})

	frontpageTemplate, _ := template.ParseFS(templates.GetFiles(), "home.html")
	err := frontpageTemplate.Execute(w, map[string]interface{}{
		"config":    viper.AllSettings(),
		"videos":    videos,
	})
	if err != nil {
		log.Fatal(err)
	}
}