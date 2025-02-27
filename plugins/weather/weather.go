package weather

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"

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
	weather := GetWeather("珠海市香洲区")
	weatherStr := fmt.Sprintf("珠海市香洲区天气: %s, 温度: %s℃, 风向: %s，风力: %s", weather.Weather, weather.Temperature, weather.WindDirection, weather.WindPower) 
	wp.WeatherChan <- weatherStr
}

func (wp *WeatherPlugin) InitHandler() {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local),
	)
	spec := "0 30 8 * * *"
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