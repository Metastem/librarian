package pages

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"codeberg.org/imabritishcow/librarian/api"
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Del("Content-Type")
		w.Header().Set("Cache-Control", "no-store")
		w.Write([]byte("500 Internal Server Error"))
		return
	}
	videos := api.GetChannelVideos(1, channel.Id)

	feed := &feeds.Feed{
		Title:       channel.Name + " - Librarian",
		Link:        &feeds.Link{Href: channel.Url},
		Image:       &feeds.Image{Url: viper.GetString("DOMAIN") + channel.Thumbnail},
		Description: channel.DescriptionTxt,
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	for i := 0; i < len(videos); i++ {
		item := &feeds.Item{
			Title:       videos[i].Title,
			Link:        &feeds.Link{Href: videos[i].Url},
			Description: videos[i].DescriptionTxt,
			Created:     time.Unix(videos[i].Timestamp, 0),
			Enclosure: 	 &feeds.Enclosure{},
		}

		if r.URL.Query().Get("enclosure") == "true" {
			url, err := url.Parse(api.GetVideoStream(videos[i].LbryUrl))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Cache-Control", "no-store")
				w.Write([]byte("500 Internal Server Error"))
				return
			}
			item.Enclosure.Url = url.String()
			item.Enclosure.Type = videos[i].MediaType
			item.Enclosure.Length = videos[i].SrcSize
		}

		feed.Items = append(feed.Items, item)
	}

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(rss))
}
