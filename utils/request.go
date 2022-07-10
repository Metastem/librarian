package utils

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var h2client = NewClient(false)
var h3client = NewClient(viper.GetBool("USE_HTTP3"))

type Data struct {
	Bytes interface{}
	Type 	string
}

func Request(url string, http3 bool, byteLimit int64, dataArr ...Data) ([]byte, error) {
	client := h2client
	if http3 {
		client = h3client
	}

	req, err := retryablehttp.NewRequest("GET", url, nil)
	data := dataArr[0]
	if data.Bytes != nil {
		req, err = retryablehttp.NewRequest("POST", url, data.Bytes)
		req.Header.Set("Content-Type", data.Type)
	}
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "https://odysee.com")
	req.Header.Set("Referer", "https://odysee.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")
	
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer res.Body.Close()

	if res.ContentLength > byteLimit {
		return []byte{}, fmt.Errorf("rejected response: over byte limit")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func RequestJSON(url string, data interface{}, http3 bool) (gjson.Result, error) {
	body, err := Request(url, http3, 1000000, Data{
		Bytes: data,
		Type: "application/json",
	})
	if err != nil {
		return gjson.Result{}, err
	}

	return gjson.Parse(string(body)), nil
}