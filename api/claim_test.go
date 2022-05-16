package api_test

import (
	"testing"

	"codeberg.org/librarian/librarian/api"
	"github.com/spf13/viper"
)

func TestGetClaim(t *testing.T) {
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	viper.Set("AUTH_TOKEN", api.NewUser())

	claim, err := api.GetClaim("@SomeOrdinaryGamers:a", "reddit-mod-gets-destroyed-on-national:b", "")
	if err != nil {
		t.Error(err)
	}
	if claim.Title != "Reddit Mod Gets Destroyed on National Television..." {
		t.Errorf("channel.Title was incorrect, got: %s, want: Reddit Mod Gets Destroyed on National Television...", claim.Title)
	}
	if claim.ClaimId != "bc3dcabb350ed498804f746c6125e1eb0127c92d" {
		t.Errorf("channel.ClaimId was incorrect, got: %s, want: bc3dcabb350ed498804f746c6125e1eb0127c92d", claim.ClaimId)
	}
	if claim.Description == "" {
		t.Fail()
	}
	if claim.Likes == 0 {
		t.Fail()
	}
	if claim.Views == 0 {
		t.Fail()
	}
}

func TestGetViews(t *testing.T) {
	viper.Set("AUTH_TOKEN", api.NewUser())

	views, err := api.GetViews("bc3dcabb350ed498804f746c6125e1eb0127c92d")
	if err != nil {
		t.Error(err)
	}
	if views == 0 {
		t.Fail()
	}
}

func TestGetLikeDislike(t *testing.T) {
	viper.Set("AUTH_TOKEN", api.NewUser())

	likesDislikes, err := api.GetLikeDislike("bc3dcabb350ed498804f746c6125e1eb0127c92d")
	if err != nil {
		t.Error(err)
	}
	if likesDislikes[0] == 0 {
		t.Fail()
	}
	if likesDislikes[1] == 0 {
		t.Fail()
	}
}