package api

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"codeberg.org/imabritishcow/librarian/types"
	"codeberg.org/imabritishcow/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var waitingComments sync.WaitGroup

func GetComments(claimId string, channelId string, channelName string) []types.Comment {
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

	commentIds := make([]string, 0)
	gjson.Get(string(commentsDataBody), "result.items.#.comment_id").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			commentIds = append(commentIds, value.String())

			return true
		},
	)

	likesDislikes := GetCommentLikeDislikes(commentIds)

	comments := make([]types.Comment, 0)

	waitingComments.Add(len(commentIds))

	gjson.Get(string(commentsDataBody), "result.items").ForEach(
		func(key, value gjson.Result) bool {
			go func() {
				timestamp := time.Unix(value.Get("timestamp").Int(), 0)

				comment := utils.ProcessText(value.Get("comment").String())

				commentId := value.Get("comment_id").String()

				comments = append(comments, types.Comment{
					Channel:   GetChannel(value.Get("channel_url").String()),
					Comment:   template.HTML(comment),
					CommentId: commentId,
					ParentId:  value.Get("parent_id").String(),
					Time:      timestamp.UTC().Format("January 2, 2006 15:04"),
					RelTime:   humanize.Time(timestamp),
					Likes:     likesDislikes[commentId][0],
					Dislikes:  likesDislikes[commentId][1],
				})

				waitingComments.Done()
			}()

			return true
		},
	)
	waitingComments.Wait()

	sort.Slice(comments[:], func(i, j int) bool {
		return comments[i].Likes < comments[j].Likes
	})

	return comments
}

func GetCommentLikeDislikes(commentIds []string) map[string][]int64 {
	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment_react_list",
		"params": map[string]interface{}{
			"comment_ids": strings.Join(commentIds, ","),
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

	likesDislikes := make(map[string][]int64)

	gjson.Get(string(commentsDataBody), "result.others_reactions").ForEach(
		func(key, value gjson.Result) bool {
			likesDislikes[key.String()] = []int64{
				value.Get("like").Int(),
				value.Get("dislike").Int(),
			}

			return true
		},
	)

	return likesDislikes
}
