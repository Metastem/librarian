package proxy

import (
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"os"
	"strings"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ProxyImage(c *fiber.Ctx) error {
	url := c.Query("url")
	hash := c.Query("hash")
	if hash == "" || url == "" {
		_, err := c.Status(400).WriteString("no hash or url")
		return err
	}

	unescapedUrl, _ := url2.QueryUnescape(url)
	unescapedUrl, _ = url2.PathUnescape(unescapedUrl)
	if !utils.VerifyHMAC(unescapedUrl, hash) {
		_, err := c.Status(400).WriteString("invalid hash")
		return err
	}

	width := c.Query("w")
	height := c.Query("h")

	optionsHash := ""
	if viper.GetString("IMAGE_CACHE") == "true" {
		hasher := sha256.New()
		hasher.Write([]byte(url + hash + width + height))
		optionsHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		image, err := os.ReadFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash)
		if err == nil {
			_, err := c.Write(image)
			return err
		}
	}

	c.Set("Cache-Control", "public,max-age=31557600")

	client := http.Client{}
	requestUrl := "https://thumbnails.odysee.com/optimize/s:" + width + ":" + height + "/quality:85/plain/" + url
	if strings.Contains(url, "static.odycdn.com/emoticons") {
		requestUrl = url
		client = *api.Client
	}
	res, err := client.Get(requestUrl)
	if err != nil {
		return err
	}

	c.Set("Content-Type", res.Header.Get("Content-Type"))

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	_, err = c.Write(data)

	if viper.GetString("IMAGE_CACHE") == "true" && res.StatusCode == 200 {
		err := os.WriteFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash, data, 0644)
		if err != nil {
			return err
		}
	}

	return err
}