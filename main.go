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

func prepareFitbit(config config.Config) fitbit.APIHandler {
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
	return fp
}

func prepareHealthplanet(config config.Config) healthplanet.APIHandler {
	hp := healthplanet.New(config.Healthplanet.ClientID, config.Healthplanet.ClientSecret)
	u := hp.OauthAuthURL()
	fmt.Printf("Visit the URL for the auth dialog: %v\n", u)

	var code string
	fmt.Print("\nplease input code:")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	hp.SetToken(code)
	return hp
}

func main() {
	config, err := config.New("")
	if err != nil {
		log.Fatal("failed to load config file", err)
	}
	_ = config
	hp := prepareHealthplanet(config)
	// prepareFitbit(config)
	data, err := hp.Innerscan("1", "6021", "6022")
	if err != nil {
		log.Fatal("failed to get innerscan data", err)
	}
	fmt.Println(data)
}
