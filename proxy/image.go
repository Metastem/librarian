package proxy

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/h2non/bimg"
	"github.com/spf13/viper"
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

	w.Header().Set("Cache-Control", "public,max-age=31557600")
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	options := bimg.Options{}
	width := r.URL.Query().Get("w")
	height := r.URL.Query().Get("h")
	crop := r.URL.Query().Get("crop")
	webp := viper.GetString("WEBP_CONVERT") == "true" && strings.Contains(r.Header.Get("Accept"), "webp")
	optionsHash := ""

	if viper.GetString("IMAGE_CACHE") == "true" {
		hasher := sha256.New()
		hasher.Write([]byte(url + hash + width + height + crop + strconv.FormatBool(webp)))
		optionsHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		image, err := os.ReadFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash)
		fmt.Println(err)
		if err == nil {
			w.Write(image)
			return
		}
	}

	if width != "" && height != "" {
		width, err := strconv.Atoi(width)
		if err != nil {
			w.Write([]byte("invalid width"))
			return
		}
		height, err := strconv.Atoi(height)
		if err != nil {
			w.Write([]byte("invalid height"))
			return
		}
		options.Width = width
		options.Height = height
	}

	if crop == "true" {
		options.Crop = true
		options.Gravity = bimg.GravityCentre
	}

	if webp {
		options.Type = bimg.WEBP
	}

	image, err := bimg.NewImage(data).Process(options)
	if err != nil {
		w.Write([]byte("error processing image"))
	}

	if viper.GetString("IMAGE_CACHE") == "true" {
		err := os.WriteFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash, image, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	w.Write(image)
}
