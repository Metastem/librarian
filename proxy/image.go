package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/h2non/bimg"
)

func ProxyImage(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	hash := r.URL.Query().Get("hash")
	if hash == "" || url == "" {
		w.WriteHeader(400)
		w.Write([]byte("no hash or url"))
		return
	}

	if !utils.VerifyHMAC(url, hash) {
		w.WriteHeader(403)
		w.Write([]byte("invalid hash"))
		return
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	options := bimg.Options{}

	if r.URL.Query().Get("w") != "" && r.URL.Query().Get("h") != "" {
		width, err := strconv.Atoi(r.URL.Query().Get("w"))
		if err != nil {
			w.Write([]byte("invalid width"))
			return
		}
		height, err := strconv.Atoi(r.URL.Query().Get("h"))
		if err != nil {
			w.Write([]byte("invalid height"))
			return
		}
		options.Width = width
		options.Height = height
	}

	if r.URL.Query().Get("crop") == "true" {
		options.Crop = true
		options.Gravity = bimg.GravityCentre
	}

	if strings.Contains(r.Header.Get("Accept"), "webp") {
		options.Type = bimg.WEBP
	}

	image, err := bimg.NewImage(data).Process(options)
	if err != nil {
		w.Write([]byte("error processing image"))
	}

	w.Header().Set("Cache-Control", "public,max-age=31557600")
	w.Write(image)
}
