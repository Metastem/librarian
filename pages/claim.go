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
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src blob: 'self'; img-src 'self'; font-src 'self'; connect-src *; media-src * blob:; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

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

	stream, err := api.GetVideoStream(claimData.LbryUrl)
	if err != nil {
		return err
	}

	switch claimData.StreamType {
	case "document":
		docRes, err := http.Get(stream)
		if err != nil {
			return err
		}

		docBody, err := ioutil.ReadAll(docRes.Body)
		if err != nil {
			return err
		}

		document := ""
		switch docRes.Header.Get("Content-Type") {
		case "text/html":
			document = utils.ProcessDocument(string(docBody), false)
		case "text/plain":
			document = string(docBody)
		case "text/markdown":
			document = utils.ProcessDocument(string(docBody), true)
		default:
			return fmt.Errorf("document type not supported: " + docRes.Header.Get("Content-Type"))
		}

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
		videoStreamType, err := api.GetVideoStreamType(stream)
		if err != nil {
			return err
		}
		isHls := false
		if videoStreamType == "application/x-mpegurl" {
			isHls = true
		}

		relatedVids, err := api.Search(claimData.Title, 1, "file", false, claimData.ClaimId, 9)
		if err != nil {
			return err
		}

		if c.Query("nojs") == "1" {
			comments := api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name, 5000, 1)

			return c.Render("claim", fiber.Map{
				"stream":         stream,
				"claim":          claimData,
				"comments":       comments,
				"commentsLength": len(comments),
				"relatedVids":    relatedVids,
				"config":         viper.AllSettings(),
				"nojs":           true,
			})
		} else {
			return c.Render("claim", fiber.Map{
				"stream":      stream,
				"streamType":  videoStreamType,
				"isHls":       isHls,
				"claim":       claimData,
				"relatedVids": relatedVids,
				"config":      viper.AllSettings(),
				"nojs":        false,
			})
		}
	case "binary":
		return c.Render("claim", fiber.Map{
			"stream":   stream,
			"download": true,
			"claim":    claimData,
			"config":   viper.AllSettings(),
		})
	default:
		if claimData.ValueType == "repost" {
			repostLink, err := utils.LbryTo(claimData.Repost)
			if err != nil {
				return err
			}
			return c.Redirect(repostLink["rel"])
		}

		if claimData.MediaType == "" && claimData.ValueType == "stream" {
			live, err := api.GetLive(claimData.Channel.Id)
			if err != nil && err.Error() != "no data associated with claim id" {
				return err
			}

			if !viper.GetBool("ENABLE_LIVE_STREAM") {
				return fmt.Errorf("live streams are disabled on this instance")
			}

			return c.Render("live", fiber.Map{
				"live":   live,
				"claim":  claimData,
				"config": viper.AllSettings(),
			})
		}

		return c.Render("claim", fiber.Map{
			"stream":   stream,
			"download": true,
			"claim":    claimData,
			"config":   viper.AllSettings(),
		})
	}
}
