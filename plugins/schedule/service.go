package schedule

import (
	"time"
)

type ScheduleService struct {
	DB ScheduleDB
}

func NewScheduleService() *ScheduleService {
	return &ScheduleService{
		DB: *NewScheduleDB(),
	}
}

func (ss *ScheduleService) AddOnce(year int, month time.Month, day, hour, minute int, content string) error {
	err := ss.DB.AddOnceTask(year, month, day, hour, minute, content)
	if err != nil {
		return err
	}
	now := time.Now()
	curYear := now.Year()
	curMonth := now.Month()
	curDay := now.Day()
	if year == curYear && month == curMonth && day == curDay {
		eventEmitter.Emit("SetTask", Date{
			Year:year,
			Month: month,
			Day: day,
			Hour: hour,
			Minute: minute,
			Content: content,
		})
	}
	return nil
}

func (ss *ScheduleService) AddLong(hour, minute int, weekday time.Weekday, content string) error {
	err := ss.DB.AddLongTask(hour, minute, weekday, content)
	if err != nil {
		return err
	}

	now := time.Now()
	
	if now.Weekday() == weekday {
		eventEmitter.Emit("SetTask", Date{
			Weekday: weekday,
			Hour: hour,
			Minute: minute,
			Content: content,
		})
	}
	return nil
}

func (ss *ScheduleService) DeleteTask(id int) error {
	err := ss.DB.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}

func (ss *ScheduleService) GetAllTasks(year int, month time.Month, day int, weekday time.Weekday) ([]Date, error) {
	tasks, err := ss.DB.GetTask(year, month, day, weekday)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
