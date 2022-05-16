package api_test

import (
	"testing"

	"codeberg.org/librarian/librarian/api"
	"github.com/spf13/viper"
)

func TestGetComments(t *testing.T) {
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	
	comments := api.GetComments("463e63afb35a319f260b36ef8d5c3dc41a98ce28", "ecf0a6be99030d0ad4e10aec11d2c0bab94246ae", "@MusicARetro", 5, 1)
	if len(comments) == 0 {
		t.Fail()
	}
	if comments[0].CommentId == "" {
		t.Fail()
	}
	if comments[0].Time == "" {
		t.Fail()
	}
}

func TestGetCommentLikeDislikes(t *testing.T) {
	likesDislikes := api.GetCommentLikeDislikes([]string{"f0023b8436278501df46ae0147cfac4a5111c05c5c9a720fd6c4a4b32227639b"})
	if len(likesDislikes) == 0 {
		t.Fail()
	}
}