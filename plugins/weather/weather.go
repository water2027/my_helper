package weather

import (
	"log"
	"os"
	"github.com/robfig/cron/v3"
	"time"

	"wx_assistant/plugins"
	"wx_assistant/utils"
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

func GetWeather(city string) Live {
	reqHelper := utils.UseRequest()
	key := os.Getenv("GAO_DE_KEY")
	var cityResponse GeocodeResponse
	err := reqHelper.Get("https://restapi.amap.com/v3/geocode/geo", utils.RequestInit{
		Query: map[string]interface{}{
			"address": city,
			"key":     key,
		},
	}, &cityResponse)
	if err != nil {
		log.Println(err)
		return Live{}
	}
	adcode := cityResponse.GeoCodes[0].AdCode
	log.Println(adcode)
	var weatherResponse WeatherResponse
	err = reqHelper.Get("https://restapi.amap.com/v3/weather/weatherInfo", 
	utils.RequestInit{
		Query: map[string]interface{}{
			"city": adcode, 
			"key": key,
			},
		}, 
	&weatherResponse)
	if err != nil {
		log.Println(err)
		return Live{}
	}
	return weatherResponse.Lives[0]
}

func (wp *WeatherPlugin) SendMessage() {
	
}

func (wp *WeatherPlugin) InitHandler() {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local),
	)
	spec := "0 0 0 * * *"
	_, err := c.AddFunc(spec, wp.SendMessage)
	if err != nil {
		panic("添加定时任务失败: " + err.Error())
	}
	c.Start()
	select {}
}

func init() {
	wp := &WeatherPlugin{
		WeatherChan: make(chan string),
	}
	plugins.RegisterPlugin(wp)
}
