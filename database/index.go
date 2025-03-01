package database

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	env := os.Getenv("GO_ENV");
	if env != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Println(err)
			return
		}
	}
	initMysqlDB()
	initRedisClient()
}