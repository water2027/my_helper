package schedule

import (
	"wx_assistant/plugins"
)

type SchedulePlugin struct {
	ScheduleChan chan string
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