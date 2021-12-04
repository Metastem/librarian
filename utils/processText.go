package utils

import (
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func ProcessText(text string, newline bool) string {
	text = string(markdown.ToHTML([]byte(text), nil, nil))
	if newline {
		text = strings.ReplaceAll(text, "\n\n", "")
		text = strings.ReplaceAll(text, "\n", "<br>")
	}
	re := regexp.MustCompile(`(?:img src=")(.*)(?:")`)
	imgs := re.FindAllString(text, len(text) / 4)
	for i := 0; i < len(imgs); i++ {
		hmac := EncodeHMAC(imgs[i])
		text = re.ReplaceAllString(text, "/image?url=$1"+hmac)
	}
	text = strings.ReplaceAll(text, `img src="`, `img src="/image?url=`)
	text = strings.ReplaceAll(text, "https://odysee.com", viper.GetString("DOMAIN"))
	text = strings.ReplaceAll(text, "https://open.lbry.com", viper.GetString("DOMAIN"))
	text = html.UnescapeString(text)
	text = bluemonday.UGCPolicy().Sanitize(text)

	return text
}

func LbryTo(link string, linkType string) string {
	link = strings.ReplaceAll(link, "#", ":")
	split := strings.Split(strings.ReplaceAll(link, "lbry://", ""), "/")
	link = "lbry://" + url.PathEscape(split[0])
	if len(split) > 1 {
		link = "lbry://" + url.PathEscape(split[0]) + "/" + url.PathEscape(split[1])
	}

	switch linkType {
	case "rel":
		link = strings.ReplaceAll(link, "lbry://", "/")
	case "http":
		link = strings.ReplaceAll(link, "lbry://", viper.GetString("DOMAIN") + "/")
	case "odysee":
		link = strings.ReplaceAll(link, "lbry://", "https://odysee.com/")
	}
	
	return link
}

func UrlEncode(link string) (string, error) {
	link2, err := url.Parse(link)
	return link2.String(), err
}