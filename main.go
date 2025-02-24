package main

import (
	"log"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"
	"wx_assistant/router"

	_ "wx_assistant/plugins/sse"
)

func main() {
	config.InitConfig()
	infoHandlers := plugins.GetHandlers()
	for _, handler := range infoHandlers {
		log.Println(handler.Name(), 1)
	}
	b := bot.NewBot(config.BotConfig.Webhook, infoHandlers)
	go func(){
		err := b.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	r := router.GetRouter()
	r.Run(":8080")
}