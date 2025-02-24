package schedule

import (
	"time"
	"wx_assistant/database"
)


type ScheduleDB struct {}

func NewScheduleDB() *ScheduleDB {
	return &ScheduleDB{}
}

func (sdb *ScheduleDB) GetTask(year int, month time.Month, day int, weekday time.Weekday) ([]Date, error) {
	var dates []Date
	db := database.GetMysqlDB()

	// 使用 OR 连接“两次查询”的条件，使用 DISTINCT 避免重复记录（根据实际情况选择是否需要 DISTINCT）
	stmt, err := db.Prepare("SELECT DISTINCT * FROM schedule WHERE (year = ? AND month = ? AND day = ?) OR weekday = ?")
	if err != nil {
		return dates, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(year, month, day, weekday)
	if err != nil {
		return dates, err
	}
	defer rows.Close()

	for rows.Next() {
		var date Date
		err := rows.Scan(&date.Id, &date.Year, &date.Month, &date.Day, &date.Weekday, &date.Hour, &date.Minute, &date.Content)
		if err != nil {
			return dates, err
		}
		dates = append(dates, date)
	}
	return dates, nil
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