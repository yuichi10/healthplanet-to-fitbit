package fitbit

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"

	"golang.org/x/oauth2"
)

type api struct {
	c     *http.Client
	oauth *oauth2.Config
	token *oauth2.Token
	code  string
}

type Fat struct {
	Date string `json:"date"`
	Time string `json:"time"`
	Fat  string `json:"fat"`
}

type Weight struct {
	Date   string
	Time   string
	Weight string
}

type APIHandler interface {
	OauthAuthURL() string
	SetToken(code string)
	Code() string
	SetBodyFats(userID string, fat ...Fat) (string, string, error)
	SetBodyWeights(userID string, weights ...Weight) (lastInsertDate, lastInsertTime string, err error)
}

func New(clientID string, clientSecret string) APIHandler {
	return &api{
		c: &http.Client{},
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"weight", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.fitbit.com/oauth2/authorize",
				TokenURL: "https://api.fitbit.com/oauth2/token",
			},
			RedirectURL: "http://localhost:8492",
		},
	}
}

func (a *api) callbackServer(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query().Get("code"))
	a.code = r.URL.Query().Get("code")
}

func (a *api) startCallbackServer() {
	http.HandleFunc("/", a.callbackServer)
	http.ListenAndServe(":8492", nil)
}
func (a *api) OauthAuthURL() string {
	go a.startCallbackServer()
	u := a.oauth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return u
}

func (a *api) Code() string {
	return a.code
}

func (a *api) SetToken(code string) {
	ctx := context.Background()

	tok, err := a.oauth.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	a.token = tok
}

func (a *api) SetBodyFats(userID string, fats ...Fat) (lastInsertDate, lastInsertTime string, err error) {
	lastInsertDate = "00000000"
	lastInsertTime = "0000"
	u, err := url.Parse("https://api.fitbit.com/1/user/")
	if err != nil {
	}
	if userID == "" {
		u.Path = path.Join(u.Path, "-", "body/log/fat.json")
	} else {
		u.Path = path.Join(u.Path, userID, "body/log/fat.json")
	}

	for _, fat := range fats {
		fmt.Println("try to set below fat data")
		fmt.Println(fat)
		u := u
		q := u.Query()
		q.Set("fat", fat.Fat)
		q.Set("date", fat.Date)
		q.Set("time", fat.Time)
		u.RawQuery = q.Encode()
		req, err := http.NewRequest("POST", u.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token.AccessToken))

		res, err := a.c.Do(req)
		if err != nil {
			fmt.Println("failed to send request to fitbit")
			return lastInsertDate, lastInsertTime, err
		}
		defer res.Body.Close()
		if res.StatusCode/100 != 2 {
			fmt.Println("failed to add data to fitbit")
			msg, _ := ioutil.ReadAll(res.Body)
			return lastInsertDate, lastInsertTime, fmt.Errorf("failed to add data to fitbit: %v", string(msg))
		}
		fmt.Println("success to set data")
		if lastInsertDate <= fat.Date {
			lastInsertDate = fat.Date
			if lastInsertTime < fat.Time {
				lastInsertTime = fat.Time
			}
		}
	}
	return lastInsertDate, lastInsertTime, err
}

func (a *api) SetBodyWeights(userID string, weights ...Weight) (lastInsertDate, lastInsertTime string, err error) {
	lastInsertDate = "00000000"
	lastInsertTime = "0000"
	u, err := url.Parse("https://api.fitbit.com/1/user/")
	if err != nil {
	}
	if userID == "" {
		u.Path = path.Join(u.Path, "-", "body/log/weight.json")
	} else {
		u.Path = path.Join(u.Path, userID, "body/log/weight.json")
	}

	for _, weight := range weights {
		fmt.Println("try to set below fat data")
		fmt.Println(weight)
		u := u
		q := u.Query()
		q.Set("weight", weight.Weight)
		q.Set("date", weight.Date)
		q.Set("time", weight.Time)
		u.RawQuery = q.Encode()
		req, err := http.NewRequest("POST", u.String(), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token.AccessToken))

		res, err := a.c.Do(req)
		if err != nil {
			fmt.Println("failed to send weight request to fitbit")
			return lastInsertDate, lastInsertTime, err
		}
		defer res.Body.Close()
		if res.StatusCode/100 != 2 {
			fmt.Println("failed to add weight data to fitbit")
			msg, _ := ioutil.ReadAll(res.Body)
			return lastInsertDate, lastInsertTime, fmt.Errorf("failed to add weight data to fitbit: %v", string(msg))
		}
		fmt.Println("success to set weight data")
		if lastInsertDate <= weight.Date {
			lastInsertDate = weight.Date
			if lastInsertTime < weight.Time {
				lastInsertTime = weight.Time
			}
		}
	}
	return
}
