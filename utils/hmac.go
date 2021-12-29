package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"

	"github.com/spf13/viper"
)

func EncodeHMAC(data string) string {
	data, _ = url.QueryUnescape(data)
	hmac := hmac.New(sha256.New, []byte(viper.GetString("HMAC_KEY")))
	hmac.Write([]byte(data))

	return hex.EncodeToString(hmac.Sum(nil))
}

func VerifyHMAC(data string, mac string) bool {
	expectedMAC := EncodeHMAC(data)
	return expectedMAC == mac
}