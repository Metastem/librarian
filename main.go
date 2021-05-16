package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imabritishcow/librarian/config"
	"github.com/imabritishcow/librarian/pages"
	"github.com/imabritishcow/librarian/templates"
)

func main() {
	config := config.GetConfig()

	fmt.Println("Librarian started on port "+config.Port)

	r := mux.NewRouter()
	r.HandleFunc("/{channel}/{video}", pages.VideoHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(templates.GetStaticFiles()))))

	http.Handle("/", r)

	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}