package proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"codeberg.org/imabritishcow/librarian/utils"
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

	w.Header().Set("Cache-Control", "public,max-age=31557600")
	w.Write(data)
}
