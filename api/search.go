package api

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/utils"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/patrickmn/go-cache"
	"github.com/tidwall/gjson"
)

var searchCache = cache.New(60*time.Minute, 30*time.Minute)

func Search(query string, page int, claimType string, nsfw bool, relatedTo string) ([]interface{}, error) {
	cacheData, found := searchCache.Get(query + fmt.Sprint(page) + claimType + fmt.Sprint(nsfw))
	if found {
		return cacheData.([]interface{}), nil
	}

	from := 0
	if page > 1 {
		from = page * 9
	}

	if len(query) <= 3 {
		return nil, fmt.Errorf("search: the query length must be between 3 and 99999")
	}

	query = strings.ReplaceAll(query, " ", "+")
	url := "https://lighthouse.odysee.com/search?s=" + query + "&free_only=true&from=" + fmt.Sprint(from) + "&nsfw=" + strconv.FormatBool(nsfw) + "&claimType=" + claimType
	if relatedTo != "" {
		url = url + "&related_to=" + relatedTo
	}

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "https://odysee.com")
	req.Header.Set("Referer", "https://odysee.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0")

	client := utils.NewClient()
	searchDataRes, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	searchDataBody, err := ioutil.ReadAll(searchDataRes.Body)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	resultsData := gjson.Parse(string(searchDataBody))

	wg := sync.WaitGroup{}
	resultsData.ForEach(
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
					results = append(results, GetChannel(value.Get("name").String()+"#"+value.Get("claimId").String(), true))
				}
			}()

			return true
		},
	)
	wg.Wait()

	searchCache.Set(query+fmt.Sprint(page)+claimType+fmt.Sprint(nsfw), results, cache.DefaultExpiration)
	return results, nil
}
