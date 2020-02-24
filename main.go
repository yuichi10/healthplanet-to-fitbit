package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yuichi10/healthplanet-to-fitbit/config"
	"github.com/yuichi10/healthplanet-to-fitbit/fitbit"
	"github.com/yuichi10/healthplanet-to-fitbit/healthplanet"
)

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

func healthTimeToFibitTime(dateTime string) (date, time string) {
	year := dateTime[0:4]
	mouth := dateTime[4:6]
	day := dateTime[6:8]
	hour := dateTime[8:10]
	min := dateTime[10:12]

	date = fmt.Sprintf("%s-%s-%s", year, mouth, day)
	time = fmt.Sprintf("%s:%s:%s", hour, min, "00")
	return
}

func FitbitTimeToHealthTime(date, time string) (dateTime string) {
	d := strings.ReplaceAll(date, "-", "")
	t := strings.ReplaceAll(time, ":", "")

	return fmt.Sprintf("%s%s", d, t)
}

func setBodyFats(config config.Config, hp healthplanet.APIHandler, fp fitbit.APIHandler) {
	data, err := hp.Innerscan("1", "6021", "6022")
	if err != nil {
		log.Fatal("failed to get innerscan data", err)
	}
	fmt.Println("row fats data")
	fmt.Println(data)
	fatsInfo := data.TagSearch("6022")
	fmt.Println("tag search fat data")
	fmt.Println(fatsInfo)

	// TODO: need to add type "0" case
	if config.LastInput.Fat.MeasureDateCase != "" {
		fmt.Println("search newer data")
		fatsInfo = fatsInfo.NewerData(config.LastInput.Fat.MeasureDateCase)
		fmt.Println(fatsInfo)
	}

	if len(fatsInfo.Data) <= 0 {
		fmt.Println("there are no fats data to send to fitbit")
		return
	}

	fats := make([]fitbit.Fat, 0, len(fatsInfo.Data))
	for _, d := range fatsInfo.Data {
		date, time := healthTimeToFibitTime(d.Date)
		f := fitbit.Fat{
			Date: date,
			Time: time,
			Fat:  d.Keydata,
		}
		fats = append(fats, f)
	}
	date, time, err := fp.SetBodyFats(config.Fitbit.UserID, fats...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added date and time")
	fmt.Println(date, time)

	t := FitbitTimeToHealthTime(date, time)
	err = config.SetLastDate("1", "6022", t)
	if err != nil {
		log.Fatal(err)
	}
}

func setBodyWeight(config config.Config, hp healthplanet.APIHandler, fp fitbit.APIHandler) {
	data, err := hp.Innerscan("1", "6021", "6022")
	if err != nil {
		log.Fatal("failed to get innerscan data", err)
	}
	weightsInfo := data.TagSearch("6021")
	fmt.Println("tag search weight data")
	fmt.Println(weightsInfo)

	// TODO: need to add type "0" case
	if config.LastInput.Weight.MeasureDateCase != "" {
		fmt.Println("search newer data")
		weightsInfo = weightsInfo.NewerData(config.LastInput.Fat.MeasureDateCase)
		fmt.Println(weightsInfo)
	}

	if len(weightsInfo.Data) <= 0 {
		fmt.Println("there are no weight data to send to fitbit")
		return
	}

	weights := make([]fitbit.Weight, 0, len(weightsInfo.Data))
	for _, d := range weightsInfo.Data {
		date, time := healthTimeToFibitTime(d.Date)
		f := fitbit.Weight{
			Date:   date,
			Time:   time,
			Weight: d.Keydata,
		}
		weights = append(weights, f)
	}
	date, time, err := fp.SetBodyWeights(config.Fitbit.UserID, weights...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("added date and time")
	fmt.Println(date, time)

	t := FitbitTimeToHealthTime(date, time)
	err = config.SetLastDate("1", "6021", t)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config, err := config.New("")
	if err != nil {
		log.Fatal("failed to load config file", err)
	}
	hp := prepareHealthplanet(config)
	fp := prepareFitbit(config)
	setBodyFats(config, hp, fp)
	setBodyWeight(config, hp, fp)
}
