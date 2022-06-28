package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"codeberg.org/librarian/librarian/data"
	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var timeRe = regexp.MustCompile(`(?m)([0-9]?[0-9]:)?[0-9]?[0-9]:[0-9]{2}`)
var ytRe = regexp.MustCompile(`https?://(www\.)?youtu\.?be(\.com)?`)
var imgurRe = regexp.MustCompile(`https?:\/\/(i\.)?imgur\.com`)
var igRe = regexp.MustCompile(`https?://(www\.)?instagram\.com`)

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
	
	replaceImgs(doc)
	replaceLinks(doc)

	text, _ = doc.Html()

	for _, match := range timeRe.FindAllString(text, -1) {
		times := strings.Split(match, ":")
		mins, _ := strconv.Atoi(times[0])
		secs, _ := strconv.Atoi(times[1])
		time := fmt.Sprint((mins * 60) + secs)
		text = strings.ReplaceAll(text, match, `<a href="#` + time + `">` + match + "</a>")
	}

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

	replaceImgs(doc)
	replaceLinks(doc)
	
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
			data.Stickers[emote] = base64.URLEncoding.EncodeToString([]byte(data.Stickers[emote]))
			proxiedImage := "/image?width=0&height=200&url=" + data.Stickers[emote] + "&hash=" + EncodeHMAC(data.Stickers[emote])
			htmlEmote := `<img loading="lazy" src="` + proxiedImage + `" height="200px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		} else if data.Emotes[emote] != "" {
			data.Emotes[emote] = base64.URLEncoding.EncodeToString([]byte(data.Emotes[emote]))
			proxiedImage := "/image?url=" + data.Emotes[emote] + "&hash=" + EncodeHMAC(data.Emotes[emote])
			htmlEmote := `<img loading="lazy" class="emote" src="` + proxiedImage + `" height="24px">`

			text = strings.ReplaceAll(text, emotes[i], htmlEmote)
		}
	}

	return text
}

func replaceImgs(doc *goquery.Document) {
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		src = base64.URLEncoding.EncodeToString([]byte(src))
		hmac := EncodeHMAC(src)
		src = "/image?url=" + src + "&hash=" + hmac
		s.SetAttr("src", src)
	})
}

func replaceLinks(doc *goquery.Document) {
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		href = strings.ReplaceAll(href, "https://odysee.com", "")
		href = strings.ReplaceAll(href, "https://lbry.tv", "")
		href = strings.ReplaceAll(href, "https://open.lbry.com", "")

		if viper.GetString("frontend.youtube") != "" {
			href = ytRe.ReplaceAllString(href, viper.GetString("frontend.youtube"))
		}
		if viper.GetString("frontend.twitter") != "" {
			href = strings.ReplaceAll(href, "https://twitter.com", viper.GetString("frontend.twitter"))
		}
		if viper.GetString("frontend.imgur") != "" {
			href = imgurRe.ReplaceAllString(href, viper.GetString("frontend.imgur"))
		}
		if viper.GetString("frontend.instagram") != "" {
			href = igRe.ReplaceAllString(href, viper.GetString("frontend.instagram"))
		}
		if viper.GetString("frontend.tiktok") != "" {
			href = strings.ReplaceAll(href, "https://tiktok.com", viper.GetString("frontend.tiktok"))
		}
		if viper.GetString("frontend.reddit") != "" {
			href = strings.ReplaceAll(href, "https://reddit.com", viper.GetString("frontend.reddit"))
		}

		s.SetAttr("href", href)
	})
}