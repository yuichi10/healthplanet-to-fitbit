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
	Healthplanet struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"healthplanet"`
	Fitbit struct {
		ClientID     string `yaml:"clientID"`
		ClientSecret string `yaml:"clientSecret"`
	} `yaml:"fitbit"`
	LastInput struct {
		// 登録日付タイプで登録されたデータをどこまで入れたか
		AddedDateCase string `yaml:"addedDateCase"`
		// 測定日付タイプで登録されたデータをどこまで入れたか
		MeasureDateCase string `yaml:"measureDateCase"`
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
	return
}
