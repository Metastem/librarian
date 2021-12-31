package pages

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ClaimHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=3600")
	c.Set("X-Frame-Options", "DENY")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=*, battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=*, geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=*, publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src blob: 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src *; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if err != nil {
		return utils.HandleError(c, err)
	}
	if claimData.ClaimId == "" {
		return c.Render("404", fiber.Map{})
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
		})
	}

	switch claimData.StreamType {
	case "document":
		docRes, err := http.Get(api.GetVideoStream(claimData.LbryUrl))
		if err != nil {
			return utils.HandleError(c, err)
		}

		if docRes.Header.Get("Content-Type") != "text/markdown" {
			return utils.HandleError(c, fmt.Errorf("document not type of text/markdown"))
		}

		docBody, err := ioutil.ReadAll(docRes.Body)
		if err != nil {
			return utils.HandleError(c, err)
		}
		document := utils.ProcessMarkdown(string(docBody))

		if c.Query("nojs") == "1" {
			comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name, 5000, 1)

			return c.Render("claim", fiber.Map{
				"document":       document,
				"claim":          claimData,
				"comments":       comments,
				"commentsLength": len(comments),
				"nojs":           true,
				"config":         viper.AllSettings(),
			})
		} else {
			return c.Render("claim", fiber.Map{
				"document": document,
				"claim":    claimData,
				"nojs":     false,
				"config":   viper.AllSettings(),
			})
		}
	case "video":
		wg := sync.WaitGroup{}

		videoStream := ""
		videoStreamType := ""
		wg.Add(1)
		go func() {
			defer wg.Done()
			videoStream = api.GetVideoStream(claimData.LbryUrl)
			videoStreamType = api.GetVideoStreamType(videoStream)
		}()

		relatedVids, err := make([]interface{}, 0), fmt.Errorf("")
		wg.Add(1)
		go func() {
			defer wg.Done()
			relatedVids, err = api.Search(claimData.Title, 1, "file", false, claimData.ClaimId)
		}()
		if err.Error() != "" {
			return utils.HandleError(c, err)
		}

		wg.Wait()
		if c.Query("nojs") == "1" {
			comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name, 5000, 1)

			return c.Render("claim", fiber.Map{
				"stream":         videoStream,
				"streamType":	 videoStreamType,
				"claim":          claimData,
				"comments":       comments,
				"commentsLength": len(comments),
				"relatedVids":    relatedVids,
				"config":         viper.AllSettings(),
				"nojs":           true,
			})
		} else {
			return c.Render("claim", fiber.Map{
				"stream":      videoStream,
				"streamType":	 videoStreamType,
				"claim":       claimData,
				"relatedVids": relatedVids,
				"config":      viper.AllSettings(),
				"nojs":        false,
			})
		}
	default:
		return utils.HandleError(c, fmt.Errorf("unsupported stream type: "+claimData.StreamType))
	}
}
