package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/types"
	"codeberg.org/librarian/librarian/utils"
	"github.com/dustin/go-humanize"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var commentCache = cache.New(30*time.Minute, 15*time.Minute)

func CommentsHandler(c *fiber.Ctx) error {
	claimId := c.Query("claim_id")
	channelId := c.Query("channel_id")
	channelName := c.Query("channel_name")
	page := c.Query("page")
	pageSize := c.Query("page_size")
	if claimId == "" || channelId == "" || channelName == "" || page == "" || pageSize == "" {
		_, err := c.Status(400).WriteString("missing query param. claim_id, channel_id, channel_name, page, page_size required")
		return err
	}

	newPage, err := strconv.Atoi(page)
	if err != nil {
		utils.HandleError(c, err)
	}
	newPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		utils.HandleError(c, err)
	}

	comments := GetComments(claimId, channelId, channelName, newPageSize, newPage)

	c.Set("Content-Type", "application/json")
	return c.JSON(map[string]interface{}{
		"comments": comments,
	})
}

func GetComments(claimId string, channelId string, channelName string, pageSize int, page int) []types.Comment {
	cacheData, found := commentCache.Get(claimId + fmt.Sprint(page) + fmt.Sprint(pageSize))
	if found {
		return cacheData.([]types.Comment)
	}

	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment.List",
		"params": map[string]interface{}{
			"page":         page,
			"claim_id":     claimId,
			"page_size":    pageSize,
			"channel_id":   channelId,
			"channel_name": channelName,
		},
	}
	commentsData, _ := json.Marshal(commentsDataMap)
	commentsDataRes, err := Client.Post("https://comments.odysee.com/api/v2?m=comment.List", "application/json", bytes.NewBuffer(commentsData))
	if err != nil {
		fmt.Println(err)
	}

	commentsDataBody, err2 := ioutil.ReadAll(commentsDataRes.Body)
	if err2 != nil {
		fmt.Println(err2)
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

	wg := sync.WaitGroup{}
	gjson.Get(string(commentsDataBody), "result.items").ForEach(
		func(key, value gjson.Result) bool {
			wg.Add(1)

			go func() {
				defer wg.Done()
				timestamp := time.Unix(value.Get("timestamp").Int(), 0)

				comment := utils.ProcessText(value.Get("comment").String(), false)

				commentId := value.Get("comment_id").String()

				comments = append(comments, types.Comment{
					Channel:   GetChannel(value.Get("channel_url").String(), false),
					Comment:   template.HTML(comment),
					CommentId: commentId,
					ParentId:  value.Get("parent_id").String(),
					Time:      timestamp.UTC().Format("January 2, 2006 15:04"),
					RelTime:   humanize.Time(timestamp),
					Likes:     likesDislikes[commentId][0],
					Dislikes:  likesDislikes[commentId][1],
				})
			}()

			return true
		},
	)
	wg.Wait()

	sort.Slice(comments[:], func(i, j int) bool {
		return comments[i].Likes > comments[j].Likes
	})

	commentCache.Set(claimId+fmt.Sprint(page)+fmt.Sprint(pageSize), comments, cache.DefaultExpiration)
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
	commentsDataRes, err := Client.Post("https://api.na-backend.odysee.com/api/v1/proxy?m=comment_react_list", "application/json", bytes.NewBuffer(commentsData))
	if err != nil {
		fmt.Println(err)
	}

	commentsDataBody, err2 := ioutil.ReadAll(commentsDataRes.Body)
	if err2 != nil {
		fmt.Println(err2)
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
