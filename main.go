package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type AccessResponse struct {
	Access_token       string
	Expires_in         int
	Expires_at         string
	Token_type         string
	Refresh_token      string
	Refresh_expires    string
	Refresh_expires_at string
	Account_id         string
	Client_id          string
	Internal_client    bool
	Client_service     string
	scope              []string
	DisplayName        string
	App                string
	In_app_id          string

	ErrorMessage string
}

type ExchangeResponse struct {
	ExpiresInSeconds int
	Code             string
	CreatingClientId string

	ErrorMessage string
}

func main() {
	var code string

	fmt.Print("Enter an auth code from https://www.epicgames.com/id/api/redirect?clientId=34a02cf8f4414e29b15921876da36f9a&responseType=code ")
	fmt.Scan(&code)

	token := authToAccess(code)
	exchange := fetchExchange(token)

	fmt.Println("Exchange Code Generated: " + exchange)
}

func authToAccess(code string) string {
	v := url.Values{}
	v.Set("grant_type", "authorization_code")
	v.Set("code", code)
	s := v.Encode()

	req, err := http.NewRequest("POST", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/token", strings.NewReader(s))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic MzRhMDJjZjhmNDQxNGUyOWIxNTkyMTg3NmRhMzZmOWE6ZGFhZmJjY2M3Mzc3NDUwMzlkZmZlNTNkOTRmYzc2Y2Y=")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	var res AccessResponse
	json.NewDecoder(resp.Body).Decode(&res)

	if res.ErrorMessage != "" {
		fmt.Println(res.ErrorMessage)
		os.Exit(0)
	}

	return res.Access_token
}

func fetchExchange(token string) string {
	req, err := http.NewRequest("GET", "https://account-public-service-prod.ol.epicgames.com/account/api/oauth/exchange", nil)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", "Bearer "+token)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	var res ExchangeResponse
	json.NewDecoder(resp.Body).Decode(&res)

	if res.ErrorMessage != "" {
		fmt.Println(res.ErrorMessage)
		os.Exit(0)
	}

	return res.Code
}
