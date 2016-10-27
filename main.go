package main

import (
	"fmt"
	"io/ioutil"
	_ "net"
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

	accessTokenCmd := kingpin.Command("access-token", "Get access token from pin")
	configFile := accessTokenCmd.Flag("config", "Configuration for client").Default("config.file").File()
	userPin := accessTokenCmd.Arg("pin", "Pin for the user account").Required().String()

	accessTokenFilename := kingpin.Flag("access-token-file", "Token containing access-token").Short('a').Default("access.token").String()

	getDaterCmd := kingpin.Command("get-dater", "Get daters")

	cmd := kingpin.Parse()

	switch cmd {
	case accessTokenCmd.FullCommand():
		var config clientConfig
		_, err := toml.DecodeReader(*configFile, &config)
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		a, err := getAccessToken(config, *userPin)
		if err != nil {
			fmt.Println("err =", err)
		} else {
			fmt.Println("got", a)
			err := ioutil.WriteFile(*accessTokenFilename, []byte(a.AccessToken), 0666)
			if err != nil {
				fmt.Println("Unable to write access token to file:", *accessTokenFilename)
			}
		}
	case getDaterCmd.FullCommand():
		accessTokenFile, err := os.Open(*accessTokenFilename)
		if err != nil {
			fmt.Println("Cannot open file:", err)
			os.Exit(1)
		}
		accessToken, err := ioutil.ReadAll(accessTokenFile)
		if err != nil {
			fmt.Println("Cannot read access token:", err)
			os.Exit(1)
		}
		err = getData(string(accessToken))
	default:
		fmt.Println("what command what")
		os.Exit(1)
	}
}
