package plugins

import (
	"log"
	"wx_assistant/bot"
)

var handlers []bot.InfoHandler

type Plugin interface {
	PluginRequired
}

type PluginRequired interface {
	bot.InfoHandler
	Name() string
	RegisterHandler() error
}

type PluginRouteOption interface {
	RegisterRoutes() error
}

func verifyRegisterRoutes(p Plugin) bool {
	_, ok := p.(PluginRouteOption)
	return ok;
}

func RegisterPlugin(p Plugin){
	err := p.RegisterHandler()
	if err != nil {
		log.Println(err)
	}
	handlers = append(handlers, p)
	if !verifyRegisterRoutes(p) {
		err = p.(PluginRouteOption).RegisterRoutes()
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func GetHandlers() []bot.InfoHandler {
	return handlers
}