package pages

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ClaimHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=3600")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=*, battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=*, geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=*, publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src blob: 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src * blob:; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if err != nil {
		return err
	}
	if claimData.ClaimId == "" {
		return c.Status(404).Render("404", fiber.Map{})
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
		})
	}

	switch claimData.StreamType {
	case "document":
		stream, err := api.GetVideoStream(claimData.LbryUrl)
		if err != nil {
			return err
		}

		docRes, err := http.Get(stream)
		if err != nil {
			return err
		}

		if docRes.Header.Get("Content-Type") != "text/markdown" {
			return fmt.Errorf("document not type of text/markdown")
		}

		docBody, err := ioutil.ReadAll(docRes.Body)
		if err != nil {
			return err
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
		videoStream, err := api.GetVideoStream(claimData.LbryUrl)
		if err != nil {
			return err
		}
		videoStreamType, err := api.GetVideoStreamType(videoStream)
		if err != nil {
			return err
		}
		isHls := false
		if videoStreamType == "application/x-mpegurl" {
			isHls = true
		}

		relatedVids, err := api.Search(claimData.Title, 1, "file", false, claimData.ClaimId)
		if err != nil {
			return err
		}

		if c.Query("nojs") == "1" {
			comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name, 5000, 1)

			return c.Render("claim", fiber.Map{
				"stream":         videoStream,
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
				"streamType":  videoStreamType,
				"isHls":       isHls,
				"claim":       claimData,
				"relatedVids": relatedVids,
				"config":      viper.AllSettings(),
				"nojs":        false,
			})
		}
	default:
		live, err := api.GetLive(claimData.Channel.Id)
		if err != nil && err.Error() != "no data associated with claim id" {
			return err
		} else if live.ClaimId != "" {
			return c.Render("live", fiber.Map{
				"live":   live,
				"claim":  claimData,
				"config": viper.AllSettings(),
			})
		}
		return fmt.Errorf("unsupported stream type: " + claimData.StreamType)
	}
}
