package pages

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/utils"
	"github.com/gorilla/feeds"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func ChannelRSSHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Cache-Control", "public,max-age=1800")
	w.Header().Set("Content-Type", "application/rss+xml")

	now := time.Now()
	channel := api.GetChannel(vars["channel"], false)
	if channel.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Del("Content-Type")
		w.Write([]byte("404 Not Found\nERROR: Unable to find channel"))
		return
	}
	claims := api.GetChannelClaims(1, channel.Id)

	image, err := utils.UrlEncode(viper.GetString("DOMAIN") + channel.Thumbnail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error\nERROR: "+err.Error()))
		return
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

		if r.URL.Query().Get("enclosure") == "true" {
			url, err := utils.UrlEncode(api.GetVideoStream(claims[i].LbryUrl))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 Internal Server Error\nERROR: "+err.Error()))
				return
			}
			item.Enclosure.Url = url
			item.Enclosure.Type = claims[i].MediaType
			item.Enclosure.Length = claims[i].SrcSize
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(rss))
}
