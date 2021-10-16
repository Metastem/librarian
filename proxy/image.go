package proxy

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"codeberg.org/imabritishcow/librarian/utils"
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

	width := r.URL.Query().Get("w")
	height := r.URL.Query().Get("h")

	optionsHash := ""
	if viper.GetString("IMAGE_CACHE") == "true" {
		hasher := sha256.New()
		hasher.Write([]byte(url + hash + width + height))
		optionsHash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

		image, err := os.ReadFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash)
		if err == nil {
			w.Write(image)
			return
		}
	}

	w.Header().Set("Cache-Control", "public,max-age=31557600")
	res, err := http.Get("https://thumbnails.odysee.com/optimize/s:" + width + ":" + height + "/quality:85/plain/" + url)
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	w.Write(data)

	if viper.GetString("IMAGE_CACHE") == "true" && res.StatusCode == 200 {
		err := os.WriteFile(viper.GetString("IMAGE_CACHE_DIR") + "/" + optionsHash, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
