package pages

import (
	"fmt"
	"sort"
	"strconv"

	"codeberg.org/librarian/librarian/api"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ChannelHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=1800")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=(), battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	theme := "light"
	if c.Cookies("theme") != "" {
		theme = c.Cookies("theme")
	}

	page := 1
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channelData, err := api.GetChannel(c.Params("channel"), true)
	if err != nil {
		return err
	}

	if channelData.Id == "" {
		return c.Status(404).Render("404", fiber.Map{})
	}

	if channelData.ValueType != "channel" {
		return ClaimHandler(c)
	}

	claims, err := api.GetChannelClaims(page, channelData.Id)
	if err != nil {
		return err
	}
	sort.Slice(claims, func(i int, j int) bool {
		return claims[i].Timestamp > claims[j].Timestamp
	})

	return c.Render("channel", fiber.Map{
		"channel":   channelData,
		"config":    viper.AllSettings(),
		"claims":    claims,
		"theme":		 theme,
		"query": map[string]interface{}{
			"page":      fmt.Sprint(page),
			"nextPage":  fmt.Sprint(page + 1),
			"prevPage":  fmt.Sprint(page - 1),
			"page0":     "0",
			"claimType": "stream",
			"stream":    "stream",
		},
	})
}
