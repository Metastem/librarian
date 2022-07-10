package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/data"
	"codeberg.org/librarian/librarian/utils"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var fpCache = cache.New(30*time.Minute, 30*time.Minute)

func GetFrontpageVideos(nsfw bool) ([]Claim, error) {
	cacheData, found := fpCache.Get("featured")
	if found {
		return cacheData.([]Claim), nil
	}

	nsfwTags := []string{"porn", "porno", "nsfw", "mature", "xxx", "sex", "creampie", "blowjob", "handjob", "vagina", "boobs", "big boobs", "big dick", "pussy", "cumshot", "anal", "hard fucking", "ass", "fuck", "hentai"}
	if nsfw {
		nsfwTags = []string{}
	}

	claimSearchData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":                12,
			"page":                     1,
			"no_totals":                true,
			"claim_type":               []string{"stream"},
			"any_tags":                 []string{},
			"not_tags":                 nsfwTags,
			"channel_ids":              data.Featured,
			"not_channel_ids":          []string{},
			"order_by":                 []string{"trending_group", "trending_mixed"},
			"fee_amount":               "<=0",
			"remove_duplicates":        true,
			"has_source":               true,
			"limit_claims_per_channel": 1,
			"release_time":             ">" + fmt.Sprint(time.Now().Unix()-15778458),
			"include_purchase_receipt": true,
		},
	}
	claimSearchReqData, _ := json.Marshal(claimSearchData)

	data, err := utils.RequestJSON(viper.GetString("API_URL")+"?m=claim_search", bytes.NewBuffer(claimSearchReqData), true)
	if err != nil {
		return []Claim{}, err
	}

	claims := make([]Claim, 0)
	wg := sync.WaitGroup{}
	data.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)
			go func() {
				defer wg.Done()

				claim, _ := ProcessClaim(value)
				claim.GetViews()
				claims = append(claims, claim)
			}()

			return true
		},
	)
	wg.Wait()

	fpCache.Set("featured", claims, cache.DefaultExpiration)
	return claims, nil
}
