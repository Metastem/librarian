package proxy

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func NewStreamProxy(addr string) {
	http.HandleFunc("/stream/", HandleStream)
	http.HandleFunc("/live/", HandleLive)

	fmt.Println("stream-proxy listening on " + addr)
	http.ListenAndServe(addr, nil)
}

func HandleStream(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(200)
		return
	}

	if r.Method != "GET" && r.Method != "HEAD" {
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
		return;
	}
 
	req, err := http.NewRequest(r.Method, "https://player.odycdn.com/" + strings.TrimPrefix(r.URL.Path, "/stream/"), nil)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return;
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
	req.Header.Set("Origin", "https://odysee.com/")
	req.Header.Set("Referer", "https://odysee.com/")
	if r.Header.Get("Range") != "" {
		req.Header.Set("Range", r.Header.Get("Range"))
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return;
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", res.Header.Get("Content-Length"))
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	if res.Header.Get("Content-Range") != "" {
		w.Header().Set("Content-Range", res.Header.Get("Content-Range"))
	}
	if res.Header.Get("Content-Type") == "application/x-mpegurl" && !strings.HasSuffix(r.URL.String(), ".m3u8") {
		w.Header().Set("Location", strings.ReplaceAll(res.Request.URL.String(), "https://player.odycdn.com", "/stream"))
		res.StatusCode = 308
	}
	w.WriteHeader(res.StatusCode)
	
	if r.Method == "GET" {
		io.Copy(w, res.Body)
		return;
	} else {
		w.Write(nil)
		return;
	}
}

func HandleLive(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(200)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
		return;
	}
 
	req, err := http.NewRequest("GET", "https://cloud.odysee.live/" + strings.TrimPrefix(r.URL.Path, "/live/"), nil)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return;
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
	req.Header.Set("Origin", "https://odysee.com/")
	req.Header.Set("Referer", "https://odysee.com/")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500 Internal Server Error"))
		return;
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.WriteHeader(res.StatusCode)

	if strings.HasSuffix(r.URL.String(), ".m3u8") {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("500 Internal Server Error"))
			return;
		}

		re := regexp.MustCompile(`(?m)^/live`)
		newBody := re.ReplaceAllString(string(body), "/live/live")
		re2 := regexp.MustCompile(`(?m)^/[0-9]{3}`)
		newBody = re2.ReplaceAllString(newBody, "/live$0")
		newBody = strings.ReplaceAll(newBody, "https://cloud.odysee.live", "/live")
		newBody = strings.ReplaceAll(newBody, "https://cdn.odysee.live", "/live")

		w.Write([]byte(newBody))
		return
	}
	
	if r.Method == "GET" {
		io.Copy(w, res.Body)
		return;
	} else {
		w.Write(nil)
		return;
	}
}