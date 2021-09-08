package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"codeberg.org/imabritishcow/librarian/api"
	"codeberg.org/imabritishcow/librarian/templates"
	"codeberg.org/imabritishcow/librarian/utils"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Cache-Control", "public,max-age=1800")

	page := 1
	pageParam, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	nsfw := false
	if r.URL.Query().Get("nsfw") == "true" {
		nsfw = true
	}

	query := r.URL.Query().Get("q")
	videoResults, err := api.Search(query, page, "file", nsfw)
	if err != nil {
		utils.HandleError(w, err)
		return
	}
	channelResults, err := api.Search(query, page, "channel", nsfw)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	searchTemplate, _ := template.ParseFS(templates.GetFiles(), "search.html")
	searchTemplate.Execute(w, map[string]interface{}{
		"videos":   videoResults,
		"channels": channelResults,
		"query": map[string]interface{}{
			"query": 		 query,
			"page":      fmt.Sprint(page),
			"nextPage":  fmt.Sprint(page + 1),
			"prevPage":  fmt.Sprint(page - 1),
		},
	})
}
