package pages

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ClaimHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=21600")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'self'; script-src blob: 'self'; connect-src *; media-src * data: blob:; block-all-mixed-content")

	theme := utils.ReadSettingFromCookie(c, "theme")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if err != nil {
		return err
	}
	if claimData.ClaimId == "" {
		return c.Status(404).Render("404", fiber.Map{"theme": theme})
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
			"theme": theme,
		})
	}

	if claimData.HasFee {
		return c.Render("errors/hasFee", fiber.Map{
			"claim": claimData,
			"theme": theme,
		})
	}

	stream, err := api.GetVideoStream(claimData.LbryUrl)
	if err != nil {
		return err
	}

	comments := []types.Comment{}
	nojs := false
	if c.Query("nojs") == "1" {
		comments = api.GetComments(claimData.ClaimId, claimData.Channel.Id, claimData.Channel.Name, 25, 1)
		nojs = true
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

		return c.Render("claim", fiber.Map{
			"document": document,
			"claim":    claimData,
			"comments": comments,
			"settings": fiber.Map{
				"theme": theme,
				"nojs":  nojs,
			},
			"config": viper.AllSettings(),
		})
	case "video":
		hls, isHls, err := api.CheckHLS(stream)
		if err != nil && err.Error() == "this content cannot be accessed due to a DMCA request" {
			return c.Status(451).Render("errors/dmca", nil)
		}
		if err != nil {
			return err
		}

		streamType := hls
		if isHls {
			streamType = "application/x-mpegurl"
			c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; img-src *; script-src blob: 'self'; connect-src *; media-src * data: blob:; block-all-mixed-content")
		}

		relatedVids, err := api.Search(claimData.Title, 1, "file", false, claimData.ClaimId, 9)
		if err != nil {
			return err
		}

		return c.Render("claim", fiber.Map{
			"stream": fiber.Map{
				"url":    stream,
				"hlsUrl": hls,
				"type":   streamType,
				"isHls":  isHls,
			},
			"comments": comments,
			"settings": fiber.Map{
				"theme": theme,
				"nojs":  nojs,
			},
			"claim":       claimData,
			"relatedVids": relatedVids,
			"config":      viper.AllSettings(),
		})
	case "binary":
		return c.Render("claim", fiber.Map{
			"stream": fiber.Map{
				"url": stream,
			},
			"download": true,
			"comments": comments,
			"claim":    claimData,
			"settings": fiber.Map{
				"theme": theme,
			},
			"config": viper.AllSettings(),
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
				return c.Render("errors/liveDisabled", fiber.Map{
					"switchUrl": c.Path(),
					"settings": fiber.Map{
						"theme": theme,
					},
				})
			}

			return c.Render("live", fiber.Map{
				"live":  live,
				"claim": claimData,
				"settings": fiber.Map{
					"theme": theme,
				},
				"config": viper.AllSettings(),
			})
		}

		return c.Render("claim", fiber.Map{
			"stream": fiber.Map{
				"url": stream,
			},
			"download": true,
			"comments": comments,
			"claim":    claimData,
			"settings": fiber.Map{
				"theme": theme,
			},
			"config": viper.AllSettings(),
		})
	}
}
