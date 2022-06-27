package utils

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

func NewClient(useHttp3 bool) *retryablehttp.Client {
	Client := retryablehttp.NewClient()
	Client.Logger = nil
	Client.RetryMax = 4
	Client.Backoff = retryablehttp.LinearJitterBackoff

	Client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
		},
		Timeout: 10 * time.Second,
	}
	if useHttp3 {
		Client.HTTPClient.Transport = &http3.RoundTripper{
			QuicConfig: &quic.Config{
				KeepAlive: true,
			},
		}
	}

	return Client
}