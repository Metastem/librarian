package api

import (
	"fmt"
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

	data, err := utils.RequestJSON(url, nil, true)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	wg := sync.WaitGroup{}
	data.ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)

			go func() {
				defer wg.Done()
				if claimType == "file" {
					vid, err := GetClaim("", value.Get("name").String(), value.Get("claimId").String())
					if err == nil && vid.ClaimId != relatedTo {
						results = append(results, vid)
					}
				} else if claimType == "channel" {
					channel, err := GetChannel(value.Get("name").String()+"#"+value.Get("claimId").String(), true)
					if err == nil {
						results = append(results, channel)
					}
				} else if claimType == "file,channel" {
					vid, err := GetClaim("", value.Get("name").String(), value.Get("claimId").String())
					if err == nil && vid.ClaimId != relatedTo {
						results = append(results, vid)
					} else if err != nil && err.Error() == "value type is channel" {
						channel, err := GetChannel(value.Get("name").String()+"#"+value.Get("claimId").String(), true)
						if err == nil {
							results = append(results, channel)
						}
					}
				}
			}()

			return true
		},
	)
	wg.Wait()

	return results, nil
}
