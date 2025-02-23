package plugins

import (
	"log"
	"github.com/gin-gonic/gin"


	"wx_assistant/router"
	"wx_assistant/bot"
)

var handlers []bot.InfoHandler

type Plugin interface {
	PluginRequired
}

type PluginRequired interface {
	bot.InfoHandler
	Name() string
}

type PluginRouteOption interface {
	RegisterRoutes(r *gin.Engine) error
}

type PluginHandlerOption interface {
	InitHandler() error
}

func verifyRegisterRoutes(p Plugin) bool {
	_, ok := p.(PluginRouteOption)
	return ok;
}

func verifyInitHandler(p Plugin) bool {
	_, ok := p.(PluginHandlerOption)
	return ok;
}

func RegisterPlugin(p Plugin){
	handlers = append(handlers, p)
	if verifyInitHandler(p) {
		err := p.(PluginHandlerOption).InitHandler()
		if err != nil {
			log.Println(err)
		}
		return
	}
	if verifyRegisterRoutes(p) {
		err := p.(PluginRouteOption).RegisterRoutes(router.GetRouter())
		if err != nil {
			log.Println(err)
		}
		return
	}
}

func GetHandlers() []bot.InfoHandler {
	return handlers
}