package main

import (
	"fmt"
	"log"

	"github.com/yuichi10/healthplanet-to-fitbit/config"
	"github.com/yuichi10/healthplanet-to-fitbit/healthplanet"
)

func main() {
	config, err := config.New("")
	if err != nil {
		log.Fatal("failed to load config file", err)
	}

	hp := healthplanet.New(config.Healthplanet.ClientID, config.Healthplanet.ClientSecret)
	u := hp.OauthAuthURL()
	fmt.Printf("Visit the URL for the auth dialog: %v\n", u)

	var code string
	fmt.Print("\nplease input code:")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	hp.SetToken(code)
}
