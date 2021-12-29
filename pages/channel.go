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
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=(), battery=(), camera=(), cross-origin-isolated=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), midi=(), navigation-override=(), payment=(), picture-in-picture=(), publickey-credentials-get=(), screen-wake-lock=(), sync-xhr=(), usb=(), web-share=(), xr-spatial-tracking=()")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	page := 1
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channelData := api.GetChannel(c.Params("channel"), true)

	if (channelData.Id == "") {
		return c.Render("404", fiber.Map{})
	}

	/*TO-DO: Add playlists

	videos := make([]types.Video, 0)
	claimType := r.URL.Query().Get("claimType")
	if claimType == "" || claimType == "stream" {
		claimType = "stream"
		videos := api.GetChannelVideos(page, channelData.Id)
		sort.Slice(videos, func (i int, j int) bool {
			return videos[i].Timestamp > videos[j].Timestamp
		})
	}*/

	claims := api.GetChannelClaims(page, channelData.Id)
	sort.Slice(claims, func(i int, j int) bool {
		return claims[i].Timestamp > claims[j].Timestamp
	})

	return c.Render("channel", fiber.Map{
		"channel":   channelData,
		"config":    viper.AllSettings(),
		"claims":    claims,
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
