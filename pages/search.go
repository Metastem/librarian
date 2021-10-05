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
			"query":    query,
			"page":     fmt.Sprint(page),
			"nextPage": fmt.Sprint(page + 1),
			"prevPage": fmt.Sprint(page - 1),
		},
	})
}
