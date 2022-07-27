package pages

import (
	"fmt"
	"strings"

	"codeberg.org/librarian/librarian/api"
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

	theme := c.Cookies("theme")
	nojs := c.Query("nojs") == "1"
	settings := fiber.Map{
		"theme": theme,
		"nojs":  nojs,
	}

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if err != nil {
		if strings.ContainsAny(err.Error(), "NOT_FOUND") {
			return c.Status(404).Render("errors/notFound", fiber.Map{"theme": theme})
		}
		return err
	}

	if claimData.ValueType == "repost" {
		repostLink, err := utils.LbryTo(claimData.Repost)
		if err != nil {
			return err
		}
		return c.Redirect(repostLink["rel"])
	}

	if utils.Contains(viper.GetStringSlice("blocked_claims"), claimData.Id) {
		return c.Status(451).Render("errors/blocked", fiber.Map{
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

	related, err := api.Search(claimData.Title, 1, "file", false, claimData.Id, 9)
	if err != nil {
		return err
	}

	if claimData.MediaType == "" && claimData.ValueType == "stream" {
		live, err := api.GetLive(claimData.Channel.Id)
		if err != nil && err.Error() != "no data associated with claim id" {
			return err
		}

		if !viper.GetBool("ENABLE_LIVE_STREAM") {
			return c.Render("errors/liveDisabled", fiber.Map{
				"switchUrl": c.Path(),
				"settings":  settings,
			})
		}

		return c.Render("live", fiber.Map{
			"live":     live,
			"claim":    claimData,
			"settings": settings,
			"config":   viper.AllSettings(),
		})
	}

	stream, err := api.GetStream(claimData.LbryUrl)
	if err != nil {
		if err.Error() == "this content cannot be accessed due to a DMCA request" {
			return c.Status(451).Render("errors/dmca", nil)
		}
		return err
	}

	comments := api.Comments{}
	if nojs {
		comments, err = claimData.GetComments("", 3, 25, 1)
		if err != nil {
			return err
		}
	}

	switch claimData.StreamType {
	case "document":
		body, err := utils.Request(stream.URL, 500000, utils.Data{Bytes: nil})
		if err != nil {
			return err
		}

		document := ""
		switch stream.Type {
		case "text/html":
			document = utils.ProcessDocument(string(body), false)
		case "text/plain":
			document = string(body)
		case "text/markdown":
			document = utils.ProcessDocument(string(body), true)
		default:
			return fmt.Errorf("document type not supported: " + stream.Type)
		}

		return c.Render("claim", fiber.Map{
			"document": document,
			"claim":    claimData,
			"comments": comments,
			"settings": settings,
			"config":   viper.AllSettings(),
		})
	case "video":
		if stream.HLS {
			c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; img-src *; script-src blob: 'self'; connect-src *; media-src * data: blob:; block-all-mixed-content")
		}

		return c.Render("claim", fiber.Map{
			"stream":      stream,
			"claim":       claimData,
			"relatedVids": related,
			"comments":    comments,
			"settings":    settings,
			"config":      viper.AllSettings(),
		})
	default:
		return c.Render("claim", fiber.Map{
			"stream":   stream,
			"download": true,
			"comments": comments,
			"claim":    claimData,
			"settings": settings,
			"config":   viper.AllSettings(),
		})
	}
}
