package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imabritishcow/librarian/pages"
	"github.com/imabritishcow/librarian/templates"
	"github.com/imabritishcow/librarian/config"
)

func main() {
	config := config.GetConfig()

	r := mux.NewRouter()
	r.HandleFunc("/{channel}/{video}", pages.VideoHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(templates.GetStaticFiles()))))

	http.Handle("/", r)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}