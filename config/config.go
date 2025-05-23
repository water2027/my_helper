package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type BotConfig struct {
	Webhook string `yaml:"webhook" json:"webhook"`
}

type AiConfig struct {
	BaseUrl    string `yaml:"base_url" json:"base_url"`
	ApiKey string `yaml:"api_key" json:"api_key"`
}

type Config struct {
	BotConfig BotConfig `yaml:"bot" json:"bot"`
	AiConfig AiConfig `yaml:"ai" json:"ai"`
}

var MyConfig Config

func InitConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	err = yaml.Unmarshal(byteValue, &MyConfig)
	if err != nil {
		log.Println(err)
		return
	}
}
