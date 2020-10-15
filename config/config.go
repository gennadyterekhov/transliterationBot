package config

import (
	"encoding/json"
	"io/ioutil"
	"transliteration_bot/helper"
)

type config struct {
	Port       int    `json:"PORT"`
	APIKey     string `json:"api_key"`
	WebhookUrl string `json:"webhook_url"`
	Cyrillic   string `json:"cyrillic"`
	Latin      string `json:"latin"`
}

func GetConfig() config {
	configStr, err := ioutil.ReadFile("config.json")
	helper.Check(err)

	var configObj config = config{}
	if err := json.Unmarshal([]byte(configStr), &configObj); err != nil {
		panic(err)
	}
	return configObj
}
