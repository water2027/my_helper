package schedule

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"

	"wx_assistant/database"
	"wx_assistant/plugins"
	"wx_assistant/utils"
)

type SchedulePlugin struct {
	ScheduleChan chan string
}

type Date struct {
	Id      int          `json:"id"`
	Year    int          `json:"year"`
	Month   time.Month   `json:"month"`
	Day     int          `json:"day"`
	Weekday time.Weekday `json:"weekday"`
	Hour    int          `json:"hour"`
	Minute  int          `json:"minute"`
	Content string       `json:"content"`
}

func (sp *SchedulePlugin) Name() string {
	return "日程安排"
}

func (sp *SchedulePlugin) SetTask(date Date) {
	go func() {
		utils.SetTodayTask(func() {
			sp.ScheduleChan <- date.Content
		}, date.Hour, date.Minute)
		database.DeleteValue(context.Background(), strconv.Itoa(date.Id))
		if date.Weekday != -1 {
			ss := NewScheduleService()
			ss.DeleteTask(date.Id)
		}
	}()
}

func (sp *SchedulePlugin) GetChan() chan string {
	return sp.ScheduleChan
}

func initTable() {
	db := database.GetMysqlDB()
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schedule (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    year INTEGER DEFAULT -1 CHECK (year >= -1),
    month INTEGER DEFAULT -1 CHECK (month BETWEEN -1 AND 12),
    day INTEGER DEFAULT -1 CHECK (day BETWEEN -1 AND 31),
    weekday INTEGER DEFAULT -1 CHECK (weekday BETWEEN -1 AND 6),
    hour INTEGER NOT NULL CHECK (hour BETWEEN 0 AND 23),
    minute INTEGER NOT NULL CHECK (minute BETWEEN 0 AND 59),
    content TEXT
)`)
	if err != nil {
		log.Println(err)
		return
	}
}

func (sp *SchedulePlugin) startDailyTask() {
	database.ClearAll(context.Background())
	for {
		select {
		case str, ok := <-sp.ScheduleChan:
			log.Println("clear", str)
			if !ok {
				return
			}
		default:
			goto ProcessTasks
		}
	}
ProcessTasks:
	// 取出今日的日程
	ss := NewScheduleService()
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	weekday := now.Weekday()
	dates, err := ss.GetAllTasks(year, month, day, weekday)
	if err != nil {
		log.Println(err)
		return
	}
	for _, date := range dates {
		database.SetValue(context.Background(), strconv.Itoa(date.Id), date.Content, time.Hour*24)
		go func() {
			utils.SetTodayTask(func() {
				sp.ScheduleChan <- date.Content
			}, date.Hour, date.Minute)
			database.DeleteValue(context.Background(), strconv.Itoa(date.Id))
			if date.Weekday != -1 {
				ss.DeleteTask(date.Id)
			}
		}()
	}
}

func (sp *SchedulePlugin) Run() {
	go func() {
		c := cron.New(
			cron.WithSeconds(),
			cron.WithLocation(time.Local),
		)
		spec := "0 0 0 * * *"
		_, err := c.AddFunc(spec, sp.startDailyTask)
		if err != nil {
			panic("添加定时任务失败: " + err.Error())
		}
		c.Start()
		select {}
	}()
}

func (sp *SchedulePlugin) InitHandler() {
	initTable()
	sp.Run()
}

var eventEmitter = *plugins.NewEventEmitter()

func init() {
	sp := &SchedulePlugin{
		ScheduleChan: make(chan string),
	}
	plugins.RegisterPlugin(sp)
	eventEmitter.On("SetTask", "SetTask", func(i ...interface{}) {
		if len(i) < 1 {
			return
		}
		date, ok := i[0].(Date)
		if ok {
			sp.SetTask(date)
		}
	})
}
