package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yuichi10/healthplanet-to-fitbit/config"
	"github.com/yuichi10/healthplanet-to-fitbit/fitbit"
	"github.com/yuichi10/healthplanet-to-fitbit/healthplanet"
)

func getFitbitCode() string {
	select {}
}

func prepareFitbit(config config.Config) {
	fp := fitbit.New(config.Fitbit.ClientID, config.Fitbit.ClientSecret)
	u := fp.OauthAuthURL()
	fmt.Printf("Visit the URL for the auth dialog: %v\n", u)

	var code string
	go func() {
		fmt.Print("\nplease input code:")
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		if code != "" {
			break
		}
		if code = fp.Code(); code != "" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	fp.SetToken(code)
}

func prepareHealthplanet(config config.Config) {
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

func main() {
	config, err := config.New("")
	if err != nil {
		log.Fatal("failed to load config file", err)
	}
	// prepareHealthplanet(config)
	prepareFitbit(config)
}
