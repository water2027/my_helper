package plugins

import (
	"log"
	"wx_assistant/message"

	"github.com/gin-gonic/gin"
)

var handlers []PluginHandlerOption

type Plugin interface {
	PluginRequired
}

type PluginRequired interface {
	Name() string
}

type PluginRouteOption interface {
	RegisterRoutes(r *gin.Engine)
}

type PluginHandlerOption interface {
	GetChan() chan message.Message
	InitHandler()
}

func verifyInitHandler(p Plugin) bool {
	_, ok := p.(PluginHandlerOption)
	return ok
}

func RegisterPlugin(p Plugin) {
	if verifyInitHandler(p) {
		handlers = append(handlers, p.(PluginHandlerOption))
		go func() {
			p.(PluginHandlerOption).InitHandler()
		}()
	}
}

func GetHandlers() []PluginHandlerOption {
	for _, h := range handlers {
		log.Println(h.(Plugin).Name(), "success")
	}
	return handlers
}
