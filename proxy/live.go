package proxy

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-retryablehttp"
)

func ProxyLive(c *fiber.Ctx) error {
	client := retryablehttp.NewClient()
	client.Logger = nil
	client.Backoff = retryablehttp.LinearJitterBackoff

	url := "https://cdn.odysee.live/" + c.Params("type") + "/" + c.Params("claimId") + "/" + c.Params("path")
	if c.Params("claimId") == "" {
		url = "https://cdn.odysee.live/" + c.Params("type") + "/" + c.Params("path")
	}
	
	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "https://odysee.com")
	req.Header.Set("Referer", "https://odysee.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	c.Set("Content-Type", res.Header.Get("Content-Type"))

	contentLen, _ := strconv.Atoi(res.Header.Get("Content-Length"))

	return c.SendStream(res.Body, contentLen)
}