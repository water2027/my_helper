package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Webhook      string `yaml:"webhook" json:"webhook"`
	TemplateStr  string `yaml:"template_str" json:"template_str"`
}

var BotConfig Config

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
	err = yaml.Unmarshal(byteValue, &BotConfig)
	if err != nil {
		log.Println(err)
		return
	}
}
