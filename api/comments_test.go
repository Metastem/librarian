package api_test

import (
	"testing"

	"codeberg.org/librarian/librarian/api"
	"codeberg.org/librarian/librarian/types"
	"github.com/spf13/viper"
)

func TestGetComments(t *testing.T) {
	viper.SetDefault("API_URL", "https://api.na-backend.odysee.com/api/v1/proxy")
	
	comments, err := api.GetComments(types.Claim{
		ClaimId: "463e63afb35a319f260b36ef8d5c3dc41a98ce28",
		Channel: types.Channel{
			Id: "ecf0a6be99030d0ad4e10aec11d2c0bab94246ae", 
			Name: "@MusicARetro", 
		},
	}, "", 3, 5, 1)
	if err != nil {
		t.Error(err)
	}
	if len(comments.Comments) == 0 {
		t.Fail()
	}
	for _, comment := range comments.Comments {
		if comment.CommentId == "" {
			t.Fail()
		}
		if comment.Time == "" {
			t.Fail()
		}
	}
}

func TestGetCommentLikeDislikes(t *testing.T) {
	likesDislikes := api.GetCommentLikeDislikes([]string{"f0023b8436278501df46ae0147cfac4a5111c05c5c9a720fd6c4a4b32227639b"})
	if len(likesDislikes) == 0 {
		t.Fail()
	}
}