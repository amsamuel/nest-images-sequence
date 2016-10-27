package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
}

func getAccessToken(c clientConfig, pin string) (accessTokenResponse, error) {
	var aResp accessTokenResponse
	resp, err := http.PostForm(NEST_API_ACCESS_URL,
		url.Values{
			"code":          {pin},
			"client_id":     {c.ClientID},
			"client_secret": {c.ClientSecret},
			"grant_type":    {"authorization_code"},
		},
	)
	if err != nil {
		return aResp, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Got Status of %d (%s): %s", resp.StatusCode, resp.Status, string(body))
		return aResp, err
	}
	if err != nil {
		return aResp, err
	}
	fmt.Println("body =", string(body))
	err = json.Unmarshal(body, &aResp)
	return aResp, err
}

func getData(accessToken string) error {

	client := &http.Client{}

	req, err := http.NewRequest("GET", NEST_API_URL, nil)
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return err
	}
	h := "Bearer " + accessToken
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", h) //"Bearer c.v8HQpTXZ84ZBx50gK3uuzrTmp3eTA7wwVbdKiviu3696s4nGqN5nP55DvAWc5AH9IomiRzdd9V94tHWTtjc3owNE3g3RH89jn9C7DiK4PWuU85yrRU0Co0FluDFphcWfyTwfsyEzPAOqVHVp")
	req.Header.Write(os.Stdout)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Got error: %s\n", err)
	} else {
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Got error: %s\n", err)
		} else {
			fmt.Printf("Got response: %s\n", string(b))
		}
	}

	return nil
}
