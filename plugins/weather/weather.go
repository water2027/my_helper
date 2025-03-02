package weather

import (
	"fmt"
	"log"
	"os"
	"time"
	"io"

	"github.com/robfig/cron/v3"

	"gopkg.in/yaml.v3"
	"wx_assistant/plugins"
	"wx_assistant/utils"
)

type WeatherPlugin struct {
	WeatherChan chan string
	WeatherConfig
}

func (wp *WeatherPlugin) Name() string {
	return "天气提醒"
}

func (wp *WeatherPlugin) GetChan() chan string {
	return wp.WeatherChan
}

func GetWeather(key, city string) Live {
	reqHelper := utils.UseRequest()
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
				"key":  key,
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
	weather := GetWeather(wp.GaoDeKey, wp.UserCity)
	weatherStr := fmt.Sprintf("%s天气: %s, 温度: %s℃, 风向: %s，风力: %s", wp.UserCity, weather.Weather, weather.Temperature, weather.WindDirection, weather.WindPower)
	wp.WeatherChan <- weatherStr
}

func (wp *WeatherPlugin) InitHandler() {
	file, err := os.Open("plugins-config/weather.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}
	var wc WeatherConfig
	err = yaml.Unmarshal(byteValue, &wc)
	if err != nil {
		log.Println(err)
		return
	}
	wp.WeatherConfig = wc
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(time.Local),
	)
	spec := "0 30 8 * * *"
	_, err = c.AddFunc(spec, wp.SendMessage)
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
