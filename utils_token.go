//Because Onedrive only use OAuth2.0

//https://docs.microsoft.com/en-us/graph/auth/auth-concepts

//we use `code flow` to get the oauth2 access_token

//https://docs.microsoft.com/en-us/advertising/guides/authentication-oauth-get-tokens?view=bingads-13
package onedrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const (
	//First, we need to access ` get_ code_ URL ` get the `code` from the attached information in the redirected URL
	get_code_url  = "https://login.live.com/oauth20_authorize.srf?response_type=code&redirect_uri=http://localhost:9999&response_mode=query&scope=offline_access%20user.read%20mail.read%20Files.ReadWrite.All&client_id=8ce1780f-12cc-45eb-a4c9-83f284e56621"
	get_token_url = "https://login.live.com/oauth20_token.srf"
)

//for test we can use the client I have registered
const (
	secret       = "-uGEL0oQZ86Wwr2I-dTF3~Y3i6sCGE5z~n"
	client_id    = "8ce1780f-12cc-45eb-a4c9-83f284e56621"
	redirect_url = "http://localhost:9999"
)

//Get some objects necessary for 'token'
type getAccessTokenObject struct {
	client_secret     string
	redirect_url      string
	client_id         string
	code              string
	has_refresh_token bool
	refresh_token     string
}

func refreshToken(o *getAccessTokenObject) (token oauth2.Token, err error) {
	if !o.has_refresh_token {
		return token, errors.New("do not have refresh token")
	}

	payload := strings.NewReader("client_id=" + o.client_id + "&scope=https%3A%2F%2Fads.microsoft.com%2Fmsads.manage&client_secret=" + o.client_secret + "&refresh_token=" + o.refresh_token + "&grant_type=refresh_token")
	client := &http.Client{}
	req, err := http.NewRequest("POST", get_token_url, payload)

	if err != nil {
		return token, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	err = json.Unmarshal(body, &token)

	if err != nil {
		return
	}

	fmt.Println(token.AccessToken, token.RefreshToken)

	return

}

func getAccessToken(o *getAccessTokenObject) (token oauth2.Token, err error) {

	payload := strings.NewReader("client_id=" + o.client_id + "&redirect_uri=" + o.redirect_url + "&client_secret=" + o.client_secret + "&code=" + o.code + "&grant_type=authorization_code")

	client := &http.Client{}
	req, err := http.NewRequest("POST", get_token_url, payload)

	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	err = json.Unmarshal(body, &token)

	if err != nil {
		return
	}

	fmt.Println(token.AccessToken, token.RefreshToken)
	o.refresh_token = token.RefreshToken
	o.has_refresh_token = true

	return

}

//获取重定向url,并没有成功获取重定向的url
func getCode(get_code_url string) (string, error) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(get_code_url)
	if err != nil {
		return "", err
	}

	url := res.Header.Get("location")
	return url, nil
}
