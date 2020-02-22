package healthplanet

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type api struct {
	c     *http.Client
	oauth *oauth2.Config
	token *oauth2.Token
}

type InnerscanData struct {
	BirthDate string `json:"birth_date"`
	Data      []struct {
		Date    string `json:"date"`
		Keydata string `json:"keydata"`
		Model   string `json:"model"`
		Tag     string `json:"tag"`
	} `json:"data"`
	Height string `json:"height"`
	Sex    string `json:"sex"`
}

type APIHandler interface {
	OauthAuthURL() string
	SetToken(code string)
	Innerscan(dateType string, tag ...string) (InnerscanData, error)
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

// Innerscan will get innerscan data from healthplanet
// dateType
// 0 : 登録日付
// 1 : 測定日付
// tag
// 6021 : 体重 (kg)
// 6022 : 体脂肪率 (%)
// 6023 : 筋肉量 (kg)
// 6024 : 筋肉スコア
// 6025 : 内臓脂肪レベル2(小数点有り、手入力含まず)
// 6026 : 内臓脂肪レベル(小数点無し、手入力含む)
// 6027 : 基礎代謝量 (kcal)
// 6028 : 体内年齢 (才)
// 6029 : 推定骨量 (kg)
func (a *api) Innerscan(dateType string, tag ...string) (InnerscanData, error) {
	u, err := url.Parse("https://www.healthplanet.jp/status/innerscan.json")
	if err != nil {
		return InnerscanData{}, nil
	}
	q := u.Query()
	q.Set("access_token", a.token.AccessToken)
	q.Set("date", dateType)
	q.Set("tag", strings.Join(tag, ","))
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	res, err := a.c.Get(u.String())
	if err != nil {
		return InnerscanData{}, nil
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		msg, _ := ioutil.ReadAll(res.Body)
		return InnerscanData{}, fmt.Errorf("innerscan return %v. %s", res.StatusCode, string(msg))
	}

	data := InnerscanData{}
	json.NewDecoder(res.Body).Decode(&data)
	return data, nil
}
