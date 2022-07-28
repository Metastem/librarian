package api_test

import (
	"testing"

	"codeberg.org/librarian/librarian/api"
	"github.com/spf13/viper"
)

func TestGetClaim(t *testing.T) {
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	viper.Set("AUTH_TOKEN", api.NewUser())

	claim, err := api.GetClaim("lbry://@SomeOrdinaryGamers:a8cca58a9a49b08a1325be5fe76646ea85201dbd/reddit-mod-gets-destroyed-on-national:bc3dcabb350ed498804f746c6125e1eb0127c92d")
	if err != nil {
		t.Error(err)
	}
	if claim.Title != "Reddit Mod Gets Destroyed on National Television..." {
		t.Errorf("claim.Title was incorrect, got: %s, want: Reddit Mod Gets Destroyed on National Television...", claim.Title)
	}
	if claim.Id != "bc3dcabb350ed498804f746c6125e1eb0127c92d" {
		t.Errorf("claim.Id was incorrect, got: %s, want: bc3dcabb350ed498804f746c6125e1eb0127c92d", claim.Id)
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

	claim := api.Claim{
		Id: "bc3dcabb350ed498804f746c6125e1eb0127c92d",
	}
	err := claim.GetViews()
	if err != nil {
		t.Error(err)
	}
	if claim.Views == 0 {
		t.Fail()
	}
}

func TestGetRatings(t *testing.T) {
	viper.Set("AUTH_TOKEN", api.NewUser())

	claim := api.Claim{
		Id: "bc3dcabb350ed498804f746c6125e1eb0127c92d",
	}
	err := claim.GetRatings()
	if err != nil {
		t.Error(err)
	}
	if claim.Likes == 0 || claim.Dislikes == 0 {
		t.Fail()
	}
}