package pages

import (
	"time"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/feeds"
	"github.com/spf13/viper"
)

func ChannelRSSHandler(c *fiber.Ctx) error {
	c.Set("Cache-Control", "public,max-age=1800")
	c.Set("Content-Type", "application/rss+xml")

	now := time.Now()
	channel := api.GetChannel(c.Params("channel"), false)
	if channel.Id == "" {
		c.Set("Content-Type", "text/plain")
		_, err := c.Status(404).WriteString("404 Not Found\nERROR: Unable to find channel")
		return err
	}
	claims := api.GetChannelClaims(1, channel.Id)

	image, err := utils.UrlEncode(viper.GetString("DOMAIN") + channel.Thumbnail)
	if err != nil {
		_, err := c.Status(500).WriteString("500 Internal Server Error\nERROR: "+err.Error())
		return err
	}

	feed := &feeds.Feed{
		Title:       channel.Name + " - Librarian",
		Link:        &feeds.Link{Href: channel.Url},
		Image:       &feeds.Image{Url: image},
		Description: channel.DescriptionTxt,
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	for i := 0; i < len(claims); i++ {
		item := &feeds.Item{
			Title:       claims[i].Title,
			Link:        &feeds.Link{Href: claims[i].Url},
			Description: "<img width=\"480\" src=\"" + viper.GetString("DOMAIN") + claims[i].ThumbnailUrl + "\"><br><br>" + string(claims[i].Description),
			Created:     time.Unix(claims[i].Timestamp, 0),
			Enclosure: 	 &feeds.Enclosure{},
		}

		if c.Query("enclosure") == "true" {
			url, err := utils.UrlEncode(api.GetVideoStream(claims[i].LbryUrl))
			if err != nil {
				_, err := c.Status(500).WriteString("500 Internal Server Error\nERROR: "+err.Error())
				return err
			}
			item.Enclosure.Url = url
			item.Enclosure.Type = claims[i].MediaType
			item.Enclosure.Length = claims[i].SrcSize
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}

	_, err = c.Write([]byte(rss))
	return err
}
