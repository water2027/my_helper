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

type Config struct {
	BotConfig BotConfig `yaml:"bot" json:"bot"`
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
