package utils

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"codeberg.org/librarian/librarian/data"
	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func ProcessText(text string, newline bool) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(text), &buf); err != nil {
		panic(err)
	}
	text = buf.String()
	if newline {
		text = strings.ReplaceAll(text, "\n\n", "")
		text = strings.ReplaceAll(text, "\n", "<br>")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		panic(err)
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		src = url.QueryEscape(src)
		hmac := EncodeHMAC(src)
		src = "/image?url=" + src + "&hash=" + hmac
		s.SetAttr("src", src)
	})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.ReplaceAll(href, "https://odysee.com", "")
		href = strings.ReplaceAll(href, "https://open.lbry.com", "")
		s.SetAttr("href", href)
	})

	text, _ = doc.Html()

	text = ReplaceStickersAndEmotes(text)

	p := bluemonday.UGCPolicy()
	p.AllowImages()
	p.RequireNoReferrerOnLinks(true)
	p.RequireNoFollowOnLinks(true)
	p.RequireCrossOriginAnonymous(true)
	text = p.Sanitize(text)

	return text
}

func ProcessDocument(text string, isMd bool) string {
	if isMd {
		md := goldmark.New(
			goldmark.WithExtensions(extension.GFM),
		)
		var buf bytes.Buffer
		if err := md.Convert([]byte(text), &buf); err != nil {
			panic(err)
		}
		text = buf.String()
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		panic(err)
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		src = url.QueryEscape(src)
		hmac := EncodeHMAC(src)
		src = "/image?url=" + src + "&hash=" + hmac
		s.SetAttr("src", src)
	})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.ReplaceAll(href, "https://odysee.com", "")
		href = strings.ReplaceAll(href, "https://open.lbry.com", "")
		s.SetAttr("href", href)
	})
	
	text, _ = doc.Html()

	p := bluemonday.UGCPolicy()
	p.AllowImages()
	p.RequireNoReferrerOnLinks(true)
	p.RequireNoFollowOnLinks(true)
	p.RequireCrossOriginAnonymous(true)
	text = p.Sanitize(text)

	return text
}

func LbryTo(link string) (map[string]string, error) {
	link = strings.ReplaceAll(link, "#", ":")
	split := strings.Split(strings.ReplaceAll(link, "lbry://", ""), "/")
	link = "lbry://" + url.PathEscape(split[0])
	if len(split) > 1 {
		link = "lbry://" + url.PathEscape(split[0]) + "/" + url.PathEscape(split[1])
	}

	link = strings.ReplaceAll(link, "lbry://", "http://domain.tld/")
	parsedLink, err := url.Parse(link)
	if err != nil {
		return map[string]string{}, err
	}
	link = parsedLink.String()

	link = strings.ReplaceAll(link, "%3A", ":")
	link = strings.ReplaceAll(link, "+", "%2B")

	return map[string]string{
		"rel":    strings.ReplaceAll(link, "http://domain.tld/", "/"),
		"http":   strings.ReplaceAll(link, "http://domain.tld/", viper.GetString("DOMAIN")+"/"),
		"odysee": strings.ReplaceAll(link, "http://domain.tld/", "https://odysee.com/"),
	}, nil
}

func UrlEncode(link string) (string, error) {
	link2, err := url.Parse(link)
	return link2.String(), err
}

func ReplaceStickersAndEmotes(text string) string {
	re := regexp.MustCompile(":(.*?):")
	emotes := re.FindAllString(text, len(text)/4)
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
