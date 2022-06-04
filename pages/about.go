package pages

import (
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
)

func AboutHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=604800")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-Robots-Tag", "noindex, noimageindex, nofollow")
	c.Set("Referrer-Policy", "no-referrer")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Strict-Transport-Security", "max-age=31557600")
	c.Set("Content-Security-Policy", "default-src 'none'; style-src 'self'; img-src 'self'; font-src 'self'; form-action 'self'; block-all-mixed-content; manifest-src 'self'")

	return c.Render("about", fiber.Map{
		"theme": utils.ReadSettingFromCookie(c, "theme"),
	})
}