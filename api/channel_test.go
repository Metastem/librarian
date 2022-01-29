package api_test

import (
	"testing"

	"codeberg.org/librarian/librarian/api"
	"github.com/spf13/viper"
)

func TestGetChannel(t *testing.T) {
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	viper.SetDefault("DOMAIN", "https://example.com")

	channel, err := api.GetChannel("@SomeOrdinaryGamers:a", true)
	if err != nil {
		t.Error(err)
	}
	if channel.Id != "a8cca58a9a49b08a1325be5fe76646ea85201dbd" {
		t.Errorf("channel.Id was incorrect, got: %s, want: a8cca58a9a49b08a1325be5fe76646ea85201dbd.", channel.Id)
	}
	if channel.Url != "https://example.com/@SomeOrdinaryGamers:a" {
		t.Errorf("channel.Id was incorrect, got: %s, want: https://example.com/@SomeOrdinaryGamers:a.", channel.Id)
	}
	if channel.RelUrl != "/@SomeOrdinaryGamers:a" {
		t.Errorf("channel.Id was incorrect, got: %s, want: /@SomeOrdinaryGamers:a.", channel.Id)
	}
	if channel.OdyseeUrl != "https://odysee.com/@SomeOrdinaryGamers:a" {
		t.Errorf("channel.Id was incorrect, got: %s, want: a8cca58a9a49b08a1325be5fe76646ea85201dbd.", channel.Id)
	}
	if channel.Thumbnail == "" {
		t.Fail()
	}
}