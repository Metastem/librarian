package proxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ProxyImage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(res.Header.Get("content-type"), "image") {
		w.Write(data)
	} else {
		w.WriteHeader(400)
	}
}