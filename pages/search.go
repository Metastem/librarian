package pages

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"codeberg.org/librarian/librarian/api"
	"github.com/gofiber/fiber/v2"
)

func SearchHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "private")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
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

	nsfw := false
	if c.Query("nsfw") == "true" {
		nsfw = true
	}

	query := c.Query("q")

	results, err := api.Search(query, page, "file,channel", nsfw, "", 12)
	if err != nil {
		return err
	}
	sort.Slice(results, func(i int, j int) bool {
		valueType := reflect.ValueOf(&results[i]).Elem().Elem().FieldByName("ValueType").String()
		if valueType == "channel" {
			return true
		} else {
			return false
		}
	})
	
	return c.Render("search", fiber.Map{
		"results":   results,
		"query": map[string]interface{}{
			"query":    query,
			"page":     fmt.Sprint(page),
			"nextPage": fmt.Sprint(page + 1),
			"prevPage": fmt.Sprint(page - 1),
		},
	})
}