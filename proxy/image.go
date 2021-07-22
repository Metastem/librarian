package proxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/h2non/bimg"
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
		switch true {
		case strings.Contains(r.Header.Get("Accept"), "image/avif"):
			data, _ = bimg.NewImage(data).Convert(bimg.AVIF)
			w.Header().Set("Content-Type", "image/avif")
		case strings.Contains(r.Header.Get("Accept"), "image/webp"):
			data, _ = bimg.NewImage(data).Convert(bimg.WEBP)
			w.Header().Set("Content-Type", "image/webp")
		default:
			data, _ = bimg.NewImage(data).Convert(bimg.PNG)
			w.Header().Set("Content-Type", "image/png")
		}

		w.Header().Set("Cache-Control", "public,max-age=31557600h")
		w.Write(data)
	} else {
		w.WriteHeader(400)
	}
}