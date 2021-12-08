package main

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/pages"
	"codeberg.org/librarian/librarian/proxy"
	"codeberg.org/librarian/librarian/templates"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/librarian/")
	viper.AddConfigPath("$HOME/.config/librarian")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "3000")
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	viper.Set("AUTH_TOKEN", api.NewUser())
	viper.WriteConfig()
	if (viper.GetString("HMAC_KEY") == "") {
		b := make([]byte, 36)
    rand.Read(b)
    viper.Set("HMAC_KEY", fmt.Sprintf("%x", b))
		viper.WriteConfig()
	}

	fmt.Println("Librarian started on port " + viper.GetString("PORT"))

	r := mux.NewRouter()
	r.HandleFunc("/", pages.FrontpageHandler)
	r.HandleFunc("/image", proxy.ProxyImage)
	r.HandleFunc("/search", pages.SearchHandler)
	r.PathPrefix("/static").Handler(http.StripPrefix("/", http.FileServer(http.FS(templates.GetStaticFiles()))))
	r.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/javascript")
		file, _ := templates.GetStaticFiles().ReadFile("static/js/sw.js")
		w.Write(file)
	})
	r.HandleFunc("/{channel}", pages.ChannelHandler)
	r.HandleFunc("/{channel}/", pages.ChannelHandler)
	r.HandleFunc("/$/invite/{channel}", pages.ChannelHandler)
	r.HandleFunc("/$/invite/{channel}/", pages.ChannelHandler)
	r.HandleFunc("/{channel}/rss", pages.ChannelRSSHandler)
	r.HandleFunc("/{channel}/{video}", pages.VideoHandler)

	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + viper.GetString("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println(srv.ListenAndServe())
}
