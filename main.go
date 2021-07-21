package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imabritishcow/librarian/pages"
	"github.com/imabritishcow/librarian/proxy"
	"github.com/imabritishcow/librarian/templates"
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
	viper.SetDefault("API_URL", "https://api.lbry.tv")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Librarian started on port " + viper.GetString("PORT"))

	r := mux.NewRouter()
	r.HandleFunc("/image", proxy.ProxyImage)
	r.PathPrefix("/static").Handler(http.StripPrefix("/", http.FileServer(http.FS(templates.GetStaticFiles()))))
	r.HandleFunc("/{channel}", pages.ChannelHandler)
	r.HandleFunc("/{channel}/{video}", pages.VideoHandler)

	http.Handle("/", r)

	err1 := http.ListenAndServe(":"+viper.GetString("PORT"), nil)
	if err1 != nil {
		log.Fatal(err)
	}
}
