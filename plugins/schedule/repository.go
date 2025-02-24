package schedule

import (
	"time"
	"wx_assistant/database"
)


type ScheduleDB struct {}

func NewScheduleDB() *ScheduleDB {
	return &ScheduleDB{}
}

func (sdb *ScheduleDB) GetTask(year, month, day, weekday int) []Date {
	var dates []Date
	db := database.GetMysqlDB()
	stmt, err := db.Prepare("select * from schedule where year = ? and month = ? and day = ?")
	if err != nil {
		return dates
	}
	rows, err := stmt.Query(year, month, day)
	if err != nil {
		return dates
	}
	defer rows.Close()
	for rows.Next() {
		var date Date
		err := rows.Scan(&date.Year, &date.Month, &date.Day, &date.Weekday, &date.Hour, &date.Minute, &date.Content)
		if err != nil {
			return dates
		}
		dates = append(dates, date)
	}
	return dates
}

func (sdb *ScheduleDB) AddOnceTask(year int, month time.Month, day, hour, minute int, content string) error {
	db := database.GetMysqlDB()
	stmt, err := db.Prepare("insert into schedule(year, month, day, hour, minute, content) values(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(year, month, day, hour, minute, content)
	if err != nil {
		return err
	}
	return nil
}

func (sdb *ScheduleDB) AddLongTask(hour, minute int, weekday time.Weekday, content string) error {
	db := database.GetMysqlDB()
	stmt, err := db.Prepare("insert into schedule(hour, minute, weekday, content) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(hour, minute, weekday, content)
	if err != nil {
		return err
	}
	return nil
}

func (sdb *ScheduleDB) DeleteTask(id int) error {
	db := database.GetMysqlDB()
	stmt, err := db.Prepare("delete from schedule where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}