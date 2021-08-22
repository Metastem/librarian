package main

import (
	"fmt"
	"net/http"
	"time"

	"codeberg.org/imabritishcow/librarian/pages"
	"codeberg.org/imabritishcow/librarian/proxy"
	"codeberg.org/imabritishcow/librarian/templates"
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

	fmt.Println("Librarian started on port " + viper.GetString("PORT"))

	r := mux.NewRouter()
	r.HandleFunc("/", pages.FrontpageHandler)
	r.HandleFunc("/image", proxy.ProxyImage)
	r.PathPrefix("/static").Handler(http.StripPrefix("/", http.FileServer(http.FS(templates.GetStaticFiles()))))
	r.HandleFunc("/{channel}", pages.ChannelHandler)
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
