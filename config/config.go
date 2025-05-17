package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type BotConfig struct {
	Name	string `yaml:"name" json:"name"`
	Webhook string `yaml:"webhook" json:"webhook"`
}

type Config struct {
	BotConfig []BotConfig `yaml:"bots" json:"bots"`
}

var MyConfig Config

func InitConfig() {
	file, err := os.Open("config.json")
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
