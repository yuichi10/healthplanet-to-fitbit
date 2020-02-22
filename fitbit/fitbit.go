package fitbit

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type api struct {
	c     *http.Client
	oauth *oauth2.Config
	token *oauth2.Token
	code  string
}

type APIHandler interface {
	OauthAuthURL() string
	SetToken(code string)
	Code() string
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
