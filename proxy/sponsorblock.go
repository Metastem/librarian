package proxy

import (
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ProxySponsorBlock(c *fiber.Ctx) error {
	url := "https://sponsor.ajay.app/api/skipSegments/" + c.Params("id")
	if c.Query("categories") == "" {
		url = "https://sponsor.ajay.app/api/skipSegments/" + c.Params("id") + "?categories=" + c.Query("categories")
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return c.Send(body)
}