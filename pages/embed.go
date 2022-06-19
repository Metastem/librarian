package pages

import (
	"fmt"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func EmbedHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=3600")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; connect-src *; media-src * blob:; block-all-mixed-content")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if claimData.ClaimId == "" {
		return c.Status(404).Render("404", fiber.Map{"theme": c.Cookies("theme")})
	}
	if err != nil {
		return err
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
			"theme": c.Cookies("theme"),
		})
	}

	if claimData.StreamType == "video" {
		videoStream, err := api.GetVideoStream(claimData.LbryUrl)
		if err != nil {
			return err
		}

		return c.Render("embed", fiber.Map{
			"stream": videoStream,
			"video":  claimData,
			"theme":  c.Cookies("theme"),
		})
	} else {
		return fmt.Errorf("unsupported stream type: " + claimData.StreamType)
	}
}
