package api

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"codeberg.org/librarian/librarian/utils"
	"github.com/tidwall/gjson"
)

func Search(query string, page int, claimType string, nsfw bool, relatedTo string, size int) ([]interface{}, error) {
	from := 0
	if page > 1 {
		from = page * size
	}

	query = strings.ReplaceAll(query, " ", "+")
	url := "https://lighthouse.odysee.tv/search?s=" + query + "&size=" + fmt.Sprint(size) + "&free_only=true&from=" + fmt.Sprint(from) + "&nsfw=" + strconv.FormatBool(nsfw) + "&claimType=" + claimType
	if relatedTo != "" {
		url = url + "&related_to=" + relatedTo
	}

	data, err := utils.RequestJSON(url, nil)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	wg := sync.WaitGroup{}

	if claimType == "channel" {
		data.ForEach(func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)

			go func() {
				defer wg.Done()
				channel, err := GetChannel("lbry://" + value.Get("name").String() + "#" + value.Get("claimId").String())
				if err == nil {
					channel.GetFollowers()
					results = append(results, channel)
				}
			}()

			return true
		})
	} else {
		urls := []string{}
		data.ForEach(func(key gjson.Result, value gjson.Result) bool {
			urls = append(urls, "lbry://" + value.Get("name").String() + "#" + value.Get("claimId").String())
			return true
		})

		claims, _ := GetClaims(urls, true, true)
		for _, claim := range claims {
			id := reflect.ValueOf(claim).FieldByName("Id").String()
			if err == nil && id != relatedTo {
				results = append(results, claim)
			}
		}
	}

	wg.Wait()

	return results, nil
}
