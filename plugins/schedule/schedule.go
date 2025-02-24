package schedule

import (
	"time"

	"wx_assistant/plugins"
)

type SchedulePlugin struct {
	ScheduleChan chan string
}

type Date struct {
	Id       int            `json:"id"`
	Year     int            `json:"year"`
	Month    time.Month     `json:"month"`
	Day      int            `json:"day"`
	Weekday  time.Weekday   `json:"weekday"`
	Weekdays []time.Weekday `json:"weekdays"`
	Hour     int            `json:"hour"`
	Minute   int            `json:"minute"`
	Content  string         `json:"content"`
}

func (sp *SchedulePlugin) Name() string {
	return "日程安排"
}

func (sp *SchedulePlugin) GetChan() chan string {
	return sp.ScheduleChan
}

func init() {
	sp := &SchedulePlugin{
		ScheduleChan: make(chan string),
	}
	plugins.RegisterPlugin(sp)
}
