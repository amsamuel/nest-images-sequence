package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "net"
	"net/http"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/alecthomas/kingpin.v2"
)

const NEST_API_URL string = "https://developer-api.nest.com"
const NEST_API_ACCESS_URL string = "https://api.home.nest.com/oauth2/access_token"

type clientConfig struct {
	ClientID     string
	ClientSecret string
}

func main() {

	var configFile = kingpin.Flag("config", "Configuration for client").Required().File()
	var userPin = kingpin.Arg("pin", "Pin for the user account").Required().String()

	kingpin.Parse()

	var config clientConfig
	_, err := toml.DecodeReader(*configFile, &config)
	if err != nil {
		fmt.Println("errorrr:", err)
		os.Exit(1)
	}
	fmt.Println("config =", config)
	fmt.Println("pin =", *userPin)

	a, err := getAccessToken(config, *userPin)
	if err != nil {
		fmt.Println("err =", err)
	} else {
		fmt.Println("got", a)
	}
}

type accessTokenResponse struct {
	AccessToken string `json:access_token`
	ExpiresIn   string `json:expires_in`
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
	err = json.Unmarshal(body, &aResp)
	return aResp, err
}
