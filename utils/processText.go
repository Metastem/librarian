package utils

import (
	"html"
	"html/template"
	"net/url"
	"regexp"
	"strings"

	"codeberg.org/librarian/librarian/data"
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
	text = bluemonday.UGCPolicy().RequireNoReferrerOnLinks(true).Sanitize(text)
	text = ReplaceStickersAndEmotes(text)

	return text
}

func ProcessMarkdown(text string) template.HTML {
	text = string(markdown.ToHTML([]byte(text), nil, nil))

	re := regexp.MustCompile(`(?:img src=")(.*)(?:")`)
	imgs := re.FindAllString(text, len(text) / 4)
	for i := 0; i < len(imgs); i++ {
		hmac := EncodeHMAC(imgs[i])
		text = re.ReplaceAllString(text, "/image?url=$1"+hmac)
	}
	text = strings.ReplaceAll(text, `img src="`, `img src="/image?url=`)

	re2 := regexp.MustCompile(`<iframe src="http(.*)>`)
	text = re2.ReplaceAllString(text, "")

	re3 := regexp.MustCompile(`<iframe src="(.*)>`)
	embeds := re3.FindAllString(text, len(text) / 4)
	for i := 0; i < len(embeds); i++ {
		embed := embeds[i]
		newEmbed := strings.ReplaceAll(embed, "#", ":")
		newEmbed = strings.ReplaceAll(newEmbed, "lbry://", "/embed/")
		text = strings.ReplaceAll(text, embed, newEmbed)
	}

	text = strings.ReplaceAll(text, "https://odysee.com", viper.GetString("DOMAIN"))
	text = strings.ReplaceAll(text, "https://open.lbry.com", viper.GetString("DOMAIN"))

	p := bluemonday.UGCPolicy()
	p.AllowElements("iframe")
	p.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	p.AllowAttrs("src").OnElements("iframe")
	text = p.Sanitize(text)
	
	return template.HTML(text)
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

func ReplaceStickersAndEmotes(text string) string {
	re := regexp.MustCompile(":(.*?):")
	emotes := re.FindAllString(text, len(text) / 4)
	for i := 0; i < len(emotes); i++ {
		emote := strings.ReplaceAll(emotes[i], ":", "")
		if data.Stickers[emote] != "" {
			proxiedImage := "/image?width=0&height=200&url=" + data.Stickers[emote] + "&hash=" + EncodeHMAC(data.Stickers[emote])
			htmlEmote := `<img loading="lazy" src="` + proxiedImage + `" height="200px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		} else if data.Emotes[emote] != "" {
			proxiedImage := "/image?url=" + data.Emotes[emote] + "&hash=" + EncodeHMAC(data.Emotes[emote])
			htmlEmote := `<img loading="lazy" class="emote" src="` + proxiedImage + `" height="24px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		}
	}

	return text
}