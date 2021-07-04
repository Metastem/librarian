package api

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"mvdan.cc/xurls/v2"
)

type Comment struct {
	Channel   Channel
	Comment   template.HTML
	CommentId string
	ParentId  string
	Time      string
	RelTime   string
	Likes     int64
	Dislikes  int64
}

func GetComments(claimId string, channelId string, channelName string) []Comment {
	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment.List",
		"params": map[string]interface{}{
			"page":         1,
			"claim_id":     claimId,
			"page_size":    99999,
			"channel_id":   channelId,
			"channel_name": channelName,
		},
	}
	commentsData, _ := json.Marshal(commentsDataMap)
	commentsDataRes, err := http.Post("https://comments.lbry.com/api/v2?m=comment.List", "application/json", bytes.NewBuffer(commentsData))
	if err != nil {
		log.Fatal(err)
	}

	commentsDataBody, err2 := ioutil.ReadAll(commentsDataRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	comments := make([]Comment, 0)

	gjson.Get(string(commentsDataBody), "result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			time := time.Unix(value.Get("timestamp").Int(), 0)

			likesDislikes := GetCommentLikeDislikes(value.Get("comment_id").String())

			comment := value.Get("comment").String()
			comment = bluemonday.UGCPolicy().Sanitize(comment)
			comment = strings.ReplaceAll(comment, "\n", "<br>")
			comment = xurls.Relaxed().ReplaceAllString(comment, "<a href=\"$1$3$4\">$1$3$4</a>")

			comments = append(comments, Comment{
				Channel:   GetChannel(value.Get("channel_url").String()),
				Comment:   template.HTML(comment),
				CommentId: value.Get("comment_id").String(),
				ParentId:  value.Get("parent_id").String(),
				Time:      time.UTC().Format("January 2, 2006 15:04"),
				RelTime:   humanize.Time(time),
				Likes:     likesDislikes[0],
				Dislikes:  likesDislikes[1],
			})

			return true
		},
	)

	return comments
}

func GetCommentLikeDislikes(commentId string) []int64 {
	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment_react_list",
		"params": map[string]interface{}{
			"comment_ids": commentId,
		},
	}
	commentsData, _ := json.Marshal(commentsDataMap)
	commentsDataRes, err := http.Post(viper.GetString("API_URL")+"/api/v1/proxy?m=comment_react_list", "application/json", bytes.NewBuffer(commentsData))
	if err != nil {
		log.Fatal(err)
	}

	commentsDataBody, err2 := ioutil.ReadAll(commentsDataRes.Body)
	if err2 != nil {
		log.Fatal(err2)
	}

	likesDislikes := gjson.Get(string(commentsDataBody), "result.others_reactions."+commentId)

	return []int64{
		likesDislikes.Get("like").Int(),
		likesDislikes.Get("dislike").Int(),
	}
}
