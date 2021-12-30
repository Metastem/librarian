package api

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/tidwall/gjson"
)

func NewUser() string {
	response, err := Client.PostForm("https://api.odysee.com/user/new", url.Values{
		"auth_token": []string{},
		"language": []string{"en"},
		"app_id": []string{"odyseecom692EAWhtoqDuAfQ6KHMXxFxt8tkhmt7sfprEMHWKjy5hf6PwZcHDV542V"},
	})
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	return gjson.Get(string(body), "data.auth_token").String()
}
