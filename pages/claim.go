package pages

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/templates"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func ClaimHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Add("Cache-Control", "public,max-age=3600")
	w.Header().Add("X-Frame-Options", "DENY")
	w.Header().Add("Referrer-Policy", "no-referrer")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("Strict-Transport-Security", "max-age=31557600")
	w.Header().Add("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=*, battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=*, geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=*, publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	w.Header().Add("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src *; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	claimData, err := api.GetClaim(vars["channel"], vars["claim"], "")
	if claimData.ClaimId == "" {
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

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		blockTemplate, _ := template.ParseFS(templates.GetFiles(), "blocked.html")
		err := blockTemplate.Execute(w, map[string]interface{}{
			"claim": claimData,
		})
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	switch claimData.StreamType {
	case "document":
		docRes, err := http.Get(api.GetVideoStream(claimData.LbryUrl))
		if err != nil {
			fmt.Println(err)
		}

		if docRes.Header.Get("Content-Type") != "text/markdown" {
			utils.HandleError(w, fmt.Errorf("document not type of text/markdown"))
		}

		docBody, err2 := ioutil.ReadAll(docRes.Body)
		if err2 != nil {
			fmt.Println(err2)
		}
		document := utils.ProcessMarkdown(string(docBody))

		comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		articleTemplate, _ := template.ParseFS(templates.GetFiles(), "article.html")
		articleTemplate.Execute(w, map[string]interface{}{
			"document":				document,
			"claim":          claimData,
			"comments":       comments,
			"commentsLength": len(comments),
			"config":         viper.AllSettings(),
		})
	case "video":
		videoStream := api.GetVideoStream(claimData.LbryUrl)
		stcStream := map[string]string{"sd": ""}
		if viper.GetString("STC_URL") != "" {
			stcStream = api.GetStcStream(claimData.ClaimId)
		}

		relatedVids, err := api.Search(claimData.Title, 1, "file", false, claimData.ClaimId)
		comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name)
		if err != nil {
			utils.HandleError(w, err)
			return
		}

		videoTemplate, _ := template.ParseFS(templates.GetFiles(), "video.html")
		videoTemplate.Execute(w, map[string]interface{}{
			"stream":         videoStream,
			"video":          claimData,
			"comments":       comments,
			"commentsLength": len(comments),
			"relatedVids":    relatedVids,
			"config":         viper.AllSettings(),
			"stcStream":      stcStream,
		})
	default:
		utils.HandleError(w, fmt.Errorf("unsupported stream type: " + claimData.StreamType))
	}
}
