package pages

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/templates"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("Cache-Control", "public,max-age=3600")
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("Referrer-Policy", "no-referrer")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("Strict-Transport-Security", "max-age=31557600")
	w.Header().Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=*, battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=*, geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=*, publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src *; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	videoData, err := api.GetVideo(vars["channel"], vars["video"], "")
	if (videoData.ClaimId == "") {
		notFoundTemplate, _ := template.ParseFS(templates.GetFiles(), "404.html")
		err := notFoundTemplate.Execute(w, nil)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	if err != nil {
		utils.HandleError(w, err)
		return
	}

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
	stcStream := map[string]string{"sd":""}
	if viper.GetString("STC_URL") != "" {
		stcStream = api.GetStcStream(videoData.ClaimId)
	}
	
	relatedVids, err := api.Search(videoData.Title, 1, "file", false, videoData.ClaimId)
	comments := api.GetComments(videoData.ClaimId, videoData.Channel.Id, videoData.Channel.Name)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	videoTemplate, _ := template.ParseFS(templates.GetFiles(), "video.html")
	videoTemplate.Execute(w, map[string]interface{}{
		"stream":         videoStream,
		"video":          videoData,
		"comments":       comments,
		"commentsLength": len(comments),
		"relatedVids":		relatedVids,
		"config":         viper.AllSettings(),
		"stcStream": 			stcStream,
	})
}
