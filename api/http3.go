package api

import (
	"net/http"
	"time"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/spf13/viper"
)

var Client = &http.Client{
	Transport: &http3.RoundTripper{},
	Timeout: 45 * time.Second,
}

func CheckUseHttp3() {
	if viper.GetBool("USE_HTTP3") {
		Client = &http.Client{}
	}
}