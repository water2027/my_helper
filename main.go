package main

import (
	"log"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"
	"wx_assistant/router"
)

func main() {
	config.InitConfig()
	b := bot.NewBot(config.BotConfig.Webhook, plugins.GetHandlers(), config.BotConfig.TemplateStr)
	err := b.Run()
	if err != nil {
		log.Println(err)
	}
	r := router.GetRouter()
	r.Run(":8080")
}