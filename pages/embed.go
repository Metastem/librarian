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
	c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=*, battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=*, geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=*, publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src * blob:; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	theme := "light"
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if claimData.ClaimId == "" {
		return c.Status(404).Render("404", fiber.Map{"theme": theme})
	}
	if err != nil {
		return err
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
			"theme": theme,
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
			"theme":  theme,
		})
	} else {
		return fmt.Errorf("unsupported stream type: " + claimData.StreamType)
	}
}
