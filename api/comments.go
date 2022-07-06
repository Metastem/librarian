package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
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

	sortBy := 3
	switch c.Query("sort_by") {
	case "controversial":
		sortBy = 2
	case "new":
		sortBy = 0
	}

	newPage, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	newPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		return err
	}

	comments, err := GetComments(types.Claim{
		ClaimId: claimId,
		Channel: types.Channel{
			Id:   channelId,
			Name: channelName,
		},
	}, c.Query("parent_id"), sortBy, newPageSize, newPage)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "application/json")
	return c.JSON(comments)
}

func GetComments(claim types.Claim, parentId string, sortBy int, pageSize int, page int) (types.Comments, error) {
	cacheData, found := commentCache.Get(claim.ClaimId + parentId + fmt.Sprint(sortBy) + fmt.Sprint(page) + fmt.Sprint(pageSize))
	if found {
		return cacheData.(types.Comments), nil
	}

	reqDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "comment.List",
		"params": map[string]interface{}{
			"page":         page,
			"claim_id":     claim.ClaimId,
			"page_size":    pageSize,
			"sort_by":      sortBy,
			"top_level":    true,
			"channel_id":   claim.Channel.Id,
			"channel_name": claim.Channel.Name,
		},
	}
	if parentId != "" {
		reqDataMap["params"].(map[string]interface{})["parent_id"] = parentId
		reqDataMap["params"].(map[string]interface{})["top_level"] = false
	}

	reqData, err := json.Marshal(reqDataMap)
	if err != nil {
		return types.Comments{}, err
	}

	data, err := utils.RequestJSON("https://comments.odysee.tv/api/v2?m=comment.List", bytes.NewBuffer(reqData), true)
	if err != nil {
		return types.Comments{}, err
	}

	commentIds := []string{}
	sortOrder := map[string]int64{}
	data.Get("result.items.#.comment_id").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			commentIds = append(commentIds, value.String())
			sortOrder[value.String()] = key.Int()
			return true
		},
	)

	likesDislikes := GetCommentLikeDislikes(commentIds)

	comments := []types.Comment{}

	wg := sync.WaitGroup{}
	data.Get("result.items").ForEach(
		func(key, value gjson.Result) bool {
			wg.Add(1)

			go func() {
				defer wg.Done()

				comment := utils.ProcessText(value.Get("comment").String(), false)
				commentId := value.Get("comment_id").String()

				timestamp := time.Unix(value.Get("timestamp").Int(), 0)
				time := timestamp.UTC().Format("January 2, 2006 15:04")
				relTime := humanize.Time(timestamp)
				if relTime == "a long while ago" {
					relTime = time
				}

				channel, err := GetChannel(value.Get("channel_url").String(), false)
				if err != nil {
					fmt.Println(err)
					return
				}

				comments = append(comments, types.Comment{
					Channel:   channel,
					Comment:   template.HTML(comment),
					CommentId: commentId,
					ParentId:  value.Get("parent_id").String(),
					Time:      time,
					RelTime:   relTime,
					Replies:   value.Get("replies").Int(),
					Likes:     likesDislikes[commentId][0],
					Dislikes:  likesDislikes[commentId][1],
				})
			}()

			return true
		},
	)
	wg.Wait()

	sort.Slice(comments, func(i, j int) bool {
		return sortOrder[comments[i].CommentId] < sortOrder[comments[j].CommentId]
	})

	returnData := types.Comments{
		Comments: comments,
		Pages:    data.Get("result.total_pages").Int(),
		Items:    data.Get("result.total_items").Int(),
	}

	commentCache.Set(claim.ClaimId+fmt.Sprint(page)+fmt.Sprint(pageSize), returnData, cache.DefaultExpiration)
	return returnData, nil
}

func GetCommentLikeDislikes(commentIds []string) map[string][]int64 {
	commentsDataMap := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "reaction.List",
		"params": map[string]interface{}{
			"comment_ids": strings.Join(commentIds, ","),
		},
	}
	commentsData, _ := json.Marshal(commentsDataMap)

	data, err := utils.RequestJSON("https://comments.odysee.tv/api/v2?m=reaction.List", bytes.NewBuffer(commentsData), true)
	if err != nil {
		fmt.Println(err)
	}

	likesDislikes := make(map[string][]int64)
	data.Get("result.others_reactions").ForEach(
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
