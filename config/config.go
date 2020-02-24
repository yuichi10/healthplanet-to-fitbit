package config

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/mitchellh/go-homedir"
)

const (
	defaultConfig = ".healthplanet_to_fitbit"
)

// Config has information to communicate healthplanet and fitbit
type Config struct {
	path         string `yaml:"-"`
	Healthplanet struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"healthplanet"`
	Fitbit struct {
		UserID       string `yaml:"userID"`
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"fitbit"`
	LastInput struct {
		Fat struct {
			// 登録日付タイプで登録されたデータをどこまで入れたか
			AddedDateCase string `yaml:"addedDateCase"`
			// 測定日付タイプで登録されたデータをどこまで入れたか
			MeasureDateCase string `yaml:"measureDateCase"`
		} `yaml:"fat"`
		Weight struct {
			// 登録日付タイプで登録されたデータをどこまで入れたか
			AddedDateCase string `yaml:"addedDateCase"`
			// 測定日付タイプで登録されたデータをどこまで入れたか
			MeasureDateCase string `yaml:"measureDateCase"`
		} `yaml:"weight"`
	} `yaml:"lastInput"`
}

// New return config structure
func New(path string) (config Config, err error) {
	if path == "" {
		home := ""
		home, err = homedir.Dir()
		if err != nil {
			return
		}
		path = filepath.Join(home, defaultConfig)
	}
	f, err := os.Open(path)
	if err != nil {
		return
	}

	d := yaml.NewDecoder(f)
	err = d.Decode(&config)
	config.path = path
	return
}

// SetLastDate set last data
// # dateType:
// 0 : 登録日付
// 1 : 測定日付
// # time:
// e.g. 202012130412 (2020年12月13日4時12分)
// # tag:
// 6021 : 体重
// 6022 : 体脂肪率
func (c *Config) SetLastDate(dateType, tag, time string) error {
	switch dateType + tag {
	case "06021":
		c.LastInput.Weight.AddedDateCase = time
	case "06022":
		c.LastInput.Fat.AddedDateCase = time
	case "16021":
		c.LastInput.Weight.MeasureDateCase = time
	case "16022":
		c.LastInput.Fat.MeasureDateCase = time
	}
	err := c.Save()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Save() error {
	f, err := os.Create(c.path)
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(f).Encode(c)
	if err != nil {
		return err
	}

	return nil
}
