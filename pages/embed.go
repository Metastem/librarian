package pages

import (
	"fmt"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
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
	c.Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self' 'unsafe-inline'; img-src 'self'; font-src 'self'; connect-src *; media-src *; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	claimData, err := api.GetClaim(c.Params("channel"), c.Params("claim"), "")
	if claimData.ClaimId == "" {
		return c.Status(404).Render("404", fiber.Map{})
	}
	if err != nil {
		return utils.HandleError(c, err)
	}

	if viper.GetString("BLOCKED_CLAIMS") != "" && strings.Contains(viper.GetString("BLOCKED_CLAIMS"), claimData.ClaimId) {
		return c.Render("blocked", fiber.Map{
			"claim": claimData,
		})
	}

	if claimData.StreamType == "video" {
		videoStream := api.GetVideoStream(claimData.LbryUrl)

		return c.Render("embed", fiber.Map{
			"stream":         videoStream,
			"video":          claimData,
		})
	} else {
		return utils.HandleError(c, fmt.Errorf("unsupported stream type: " + claimData.StreamType))
	}
}
