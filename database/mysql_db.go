package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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

func initMysqlDB() error {
	var err error
	once.Do(func() {
		mysqlHost := os.Getenv("MYSQL_HOST")
		mysqlUser := os.Getenv("MYSQL_USER")
		mysqlPassword := os.Getenv("MYSQL_PASSWORD")
		mysqlPort := os.Getenv("MYSQL_PORT")
		mysqlDB := os.Getenv("MYSQL_DB")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)
		log.Println(dsn)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Println(err)
			return
		}

		err = db.Ping()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("成功连接到MySQL数据库")
	})
	return err
}
