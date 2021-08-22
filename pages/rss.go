package pages

import (
	"fmt"
	"net/http"
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
	channel := api.GetChannel(vars["channel"])
	videos := api.GetChannelVideos(1, channel.Id)

	feed := &feeds.Feed{
		Title:       channel.Name + " - Librarian",
		Link:        &feeds.Link{Href: channel.Url},
		Image:       &feeds.Image{Url: "https://" + viper.GetString("DOMAIN") + channel.Thumbnail},
		Description: channel.DescriptionTxt,
		Created:     now,
	}

	feed.Items = []*feeds.Item{}

	for i := 0; i < len(videos); i++ {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       videos[i].Title,
			Link:        &feeds.Link{Href: videos[i].Url},
			Description: videos[i].DescriptionTxt,
			Created:     time.Unix(videos[i].Timestamp, 0),
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		fmt.Println(err)
	}

	w.Write([]byte(rss))
}
