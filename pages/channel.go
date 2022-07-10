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
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	page := 1
	pageParam, err := strconv.Atoi(c.Query("page"))
	if err == nil || pageParam != 0 {
		page = pageParam
	}

	channel, err := api.GetChannel(c.Params("channel"))
	if err != nil {
		return err
	}
	channel.GetFollowers()

	if channel.Id == "" {
		return c.Status(404).Render("errors/notFound", fiber.Map{})
	}

	if channel.ValueType != "channel" {
		return ClaimHandler(c)
	}

	claims, err := channel.GetClaims(page)
	if err != nil {
		return err
	}
	sort.Slice(claims, func(i int, j int) bool {
		return claims[i].Timestamp > claims[j].Timestamp
	})

	return c.Render("channel", fiber.Map{
		"channel": channel,
		"config":  viper.AllSettings(),
		"claims":  claims,
		"theme":   c.Cookies("theme"),
		"query": fiber.Map{
			"page":        fmt.Sprint(page),
			"prevPageIs0": (page - 1) == 0,
			"nextPage":    fmt.Sprint(page + 1),
			"prevPage":    fmt.Sprint(page - 1),
		},
	})
}
