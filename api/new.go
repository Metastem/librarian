package api

import (
	"io/ioutil"
	"log"
	"net/url"

	"codeberg.org/librarian/librarian/utils"
	"github.com/tidwall/gjson"
)

func NewUser() string {
	Client := utils.NewClient()
	response, err := Client.PostForm("https://api.odysee.com/user/new", url.Values{
		"auth_token": []string{},
		"language": []string{"en"},
		"app_id": []string{"odyseecom692EAWhtoqDuAfQ6KHMXxFxt8tkhmt7sfprEMHWKjy5hf6PwZcHDV542V"},
	})
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return gjson.Get(string(body), "data.auth_token").String()
}
