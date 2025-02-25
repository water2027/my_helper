package main

import (
	"log"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"
	"wx_assistant/router"
	
	_ "wx_assistant/database"
	_ "wx_assistant/plugins/schedule"
	_ "wx_assistant/plugins/sse"
)

func main() {
	config.InitConfig()
	r := router.GetRouter()
	r.LoadHTMLGlob("templates/**/*")
	infoHandlers := plugins.GetHandlers()
	b := bot.NewBot(config.BotConfig.Webhook, infoHandlers)
	go func() {
		err := b.Run()
		if err != nil {
			log.Println(err)
		}
	}()
	r.Run(":8080")
}
