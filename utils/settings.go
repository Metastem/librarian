package utils

import "github.com/gofiber/fiber/v2"

func ReadSettingFromCookie(c *fiber.Ctx, name string) string {
	if c.Cookies(name) != "" {
		return c.Cookies(name)
	} else {
		return "false"
	}
}