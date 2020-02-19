package healthplanet

import (
	"context"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

type api struct {
	c     *http.Client
	oauth *oauth2.Config
	token *oauth2.Token
}

type APIHandler interface {
	OauthAuthURL() string
	SetToken(code string)
}

func New(clientID string, clientSecret string) APIHandler {
	return &api{
		c: &http.Client{},
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"innerscan", "sphygmomanometer", "pedometer", "smug"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.healthplanet.jp/oauth/auth",
				TokenURL: "https://www.healthplanet.jp/oauth/token",
			},
			RedirectURL: "https://www.healthplanet.jp/success.html",
		},
	}
}

func (a *api) OauthAuthURL() string {
	row := a.oauth.AuthCodeURL("state", oauth2.AccessTypeOffline)
	u := strings.ReplaceAll(row, "+", ",")

	return u
}

func (a *api) SetToken(code string) {
	ctx := context.Background()

	tok, err := a.oauth.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	a.token = tok
}
