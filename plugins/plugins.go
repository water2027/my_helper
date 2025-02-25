package plugins

import (
	"github.com/gin-gonic/gin"

	"wx_assistant/router"
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
	GetChan() chan string
	InitHandler()
}

func verifyRegisterRoutes(p Plugin) bool {
	_, ok := p.(PluginRouteOption)
	return ok
}

func verifyInitHandler(p Plugin) bool {
	_, ok := p.(PluginHandlerOption)
	return ok
}

func RegisterPlugin(p Plugin) {
	if verifyInitHandler(p) {
		go func() {
			p.(PluginHandlerOption).InitHandler()
			handlers = append(handlers, p.(PluginHandlerOption))
		}()
	}
	if verifyRegisterRoutes(p) {
		go func() {
			p.(PluginRouteOption).RegisterRoutes(router.GetRouter())
		}()
	}
}

func GetHandlers() []PluginHandlerOption {
	return handlers
}
