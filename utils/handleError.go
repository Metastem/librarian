package utils

import (
	"github.com/gofiber/fiber/v2"
)

func HandleError(c *fiber.Ctx, err error) error {
	c.Status(500)
	return c.Render("error", fiber.Map{
		"err": err,
	})
}