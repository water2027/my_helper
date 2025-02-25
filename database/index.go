package database

import (
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
		return
	}
	initMysqlDB()
	initRedisClient()
}