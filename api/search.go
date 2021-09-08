package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

var waitingResults sync.WaitGroup

func Search(query string, page int, claimType string, nsfw bool) ([]interface{}, error) {
	from := 0
	if page > 1 {
		from = page * 9
	}

	if len(query) <= 3 {
		return nil, fmt.Errorf("search: the query length must be between 3 and 99999")
	}

	query = strings.ReplaceAll(query, " ", "+")
	searchDataRes, err := http.Get("https://lighthouse.lbry.com/search?s=" + query + "&free_only=true&from=" + fmt.Sprint(from) + "&nsfw=" + strconv.FormatBool(nsfw) + "&claimType=" + claimType)
	if err != nil {
		return nil, err
	}

	searchDataBody, err := ioutil.ReadAll(searchDataRes.Body)
	if err != nil {
		return nil, err
	}

	results := make([]interface{}, 0)
	resultsData := gjson.Parse(string(searchDataBody))

	waitingResults.Add(int(resultsData.Get("#").Int()))
	resultsData.ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			go func() {
				if (claimType == "file") {
					results = append(results, GetVideo("", value.Get("name").String(), value.Get("claimId").String()))
				} else if (claimType == "channel") {
					results = append(results, GetChannel(value.Get("name").String() + "#" + value.Get("claimId").String(), true))
				}

				waitingResults.Done()
			}()

			return true
		},
	)
	waitingResults.Wait()

	return results, nil
}