package api

import (
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/spf13/viper"
)

var Client = &http.Client{
	Transport: &http3.RoundTripper{},
}

func CheckUseHttp3() {
	if viper.GetBool("USE_HTTP3") == true {
		Client = &http.Client{}
	}
}