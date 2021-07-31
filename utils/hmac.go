package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/spf13/viper"
)

func EncodeHMAC(data string) string {
	hmac := hmac.New(sha256.New, []byte(viper.GetString("HMAC_KEY")))
	hmac.Write([]byte(data))

	return hex.EncodeToString(hmac.Sum(nil))
}

func VerifyHMAC(data string, mac string) bool {
	h := hmac.New(sha256.New, []byte(viper.GetString("HMAC_KEY")))
	h.Write([]byte(data))
	expectedMAC := hex.EncodeToString(h.Sum(nil))

	return expectedMAC == mac
}