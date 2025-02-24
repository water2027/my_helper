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
	return nil
}

func (ss *ScheduleService) AddLong(hour, minute int, weekday time.Weekday, content string) error {
	err := ss.DB.AddLongTask(hour, minute, weekday, content)
	if err != nil {
		return err
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

func (ss *ScheduleService) AddPage() {
	
}

func (ss *ScheduleService) BrowsePage() {}
