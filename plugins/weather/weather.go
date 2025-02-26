package weather

import (
	"wx_assistant/plugins"
)

type WeatherPlugin struct {
	WeatherChan chan string
}

func (wp *WeatherPlugin) Name() string {
	return "天气提醒"
}

func (wp *WeatherPlugin) GetChan() chan string {
	return wp.WeatherChan
}



func init() {
	wp := &WeatherPlugin{
		WeatherChan: make(chan string),
	}
	plugins.RegisterPlugin(wp)
}