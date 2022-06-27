package proxy

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
)

func ProxyLive(c *fiber.Ctx) error {
	client := utils.NewClient(false)

	url := "https://cloud.odysee.live/" + c.Params("+")

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

	if strings.Contains(res.Header.Get("Content-Type"), "text/plain") {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		newBody := strings.ReplaceAll(string(body), "https://cloud.odysee.live", "/live")
		newBody = strings.ReplaceAll(newBody, "https://cdn.odysee.live", "/live")
		re := regexp.MustCompile(`(?m)^/[0-9]{3}`)
		newBody = re.ReplaceAllString(newBody, "/live$0")
		if viper.GetString("LIVE_STREAMING_URL") != "" {
			newBody = strings.ReplaceAll(string(body), "https://cloud.odysee.live", viper.GetString("LIVE_STREAMING_URL"))
			newBody = strings.ReplaceAll(newBody, "https://cdn.odysee.live", viper.GetString("LIVE_STREAMING_URL"))
			newBody = re.ReplaceAllString(newBody, viper.GetString("LIVE_STREAMING_URL") + "$0")
		}
		return c.Send([]byte(newBody))
	}

	return c.SendStream(res.Body, contentLen)
}