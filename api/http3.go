package api

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/spf13/viper"
)

var Client = retryablehttp.NewClient()

func CheckUseHttp3() {
	Client.Logger = nil
	Client.RetryMax = 4

	Client.HTTPClient = &http.Client{}
	if viper.GetBool("USE_HTTP3") {
		Client.HTTPClient = &http.Client{
			Transport: &http3.RoundTripper{},
		}
	}
}