package database

import (
	"database/sql"
	"fmt"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
	mu   sync.Mutex
)

func GetMysqlDB() *sql.DB {
	mu.Lock()
	defer mu.Unlock()
	return db
}

func initTable() error {
	return nil
}

func initMysqlDB() error {
	var err error
	once.Do(func() {
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mysqlPort := os.Getenv("MYSQL_PORT")
		mysqlDB := os.Getenv("MYSQL_DB")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			return
		}

		err = db.Ping()
		if err != nil {
			return
		}
		err = initTable()
		if err != nil {
			return
		}
		fmt.Println("成功连接到MySQL数据库")
	})
	return err
}
