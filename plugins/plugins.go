package plugins

import (
	"log"
	"github.com/gin-gonic/gin"


	"wx_assistant/router"
)

var handlers []Plugin

type Plugin interface {
	PluginRequired
}

type PluginRequired interface {
	GetChan() chan string
	Name() string
}

type PluginRouteOption interface {
	RegisterRoutes(r *gin.Engine)
}

type PluginHandlerOption interface {
	InitHandler()
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
		log.Println("init handler")
		go func(){
			p.(PluginHandlerOption).InitHandler()
		}()
	}
	if verifyRegisterRoutes(p) {
		go func(){
			p.(PluginRouteOption).RegisterRoutes(router.GetRouter())
		}()
	}
}

func GetHandlers() []Plugin {
	return handlers
}