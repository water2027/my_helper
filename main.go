package main

import (
	"log"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"
	"wx_assistant/router"

	_ "wx_assistant/plugins/sse"
	_ "wx_assistant/plugins/schedule"
)

func main() {
	config.InitConfig()
	r := router.GetRouter()
	r.LoadHTMLGlob("templates/**/*")
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
	r.Run(":8080")
}