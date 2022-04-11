package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"codeberg.org/librarian/librarian/data"
	"codeberg.org/librarian/librarian/types"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var fpCache = cache.New(30*time.Minute, 30*time.Minute)

func GetFrontpageVideos() ([]types.Claim, error) {
	cacheData, found := fpCache.Get("featured")
	if found {
		return cacheData.([]types.Claim), nil
	}

	claimSearchData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "claim_search",
		"params": map[string]interface{}{
			"page_size":                12,
			"page":											1,
			"no_totals":                true,
			"claim_type":               []string{"stream"},
			"any_tags":                 []string{},
			"not_tags":                 []string{"porn", "porno", "nsfw", "mature", "xxx", "sex", "creampie", "blowjob", "handjob", "vagina", "boobs", "big boobs", "big dick", "pussy", "cumshot", "anal", "hard fucking", "ass", "fuck", "hentai"},
			"channel_ids":              data.Featured,
			"not_channel_ids":          []string{},
			"order_by":                 []string{"trending_group", "trending_mixed"},
			"fee_amount":               "<=0",
			"remove_duplicates":				true,
			"has_source":								true,
			"limit_claims_per_channel":	1,
			"release_time":             ">" + fmt.Sprint(time.Now().Unix()-15778458),
			"include_purchase_receipt": true,
		},
	}
	claimSearchReqData, _ := json.Marshal(claimSearchData)
	frontpageDataRes, err := http.Post(viper.GetString("API_URL")+"?m=claim_search", "application/json", bytes.NewBuffer(claimSearchReqData))
	if err != nil {
		return []types.Claim{}, err
	}

	frontpageDataBody, err := ioutil.ReadAll(frontpageDataRes.Body)
	if err != nil {
		return []types.Claim{}, err
	}

	claims := make([]types.Claim, 0)
	claimsData := gjson.Parse(string(frontpageDataBody))

	wg := sync.WaitGroup{}
	claimsData.Get("result.items").ForEach(
		func(key gjson.Result, value gjson.Result) bool {
			wg.Add(1)
			go func() {
				defer wg.Done()

				claim, _ := ProcessClaim(value, true, false)
				claims = append(claims, claim)
			}()

			return true
		},
	)
	wg.Wait()

	fpCache.Set("featured", claims, cache.DefaultExpiration)
	return claims, nil
}
