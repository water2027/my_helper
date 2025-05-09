package init

import (
	// 读取环境变量
	"github.com/joho/godotenv"

	"wx_assistant/config"
)

func init() {
	// 读取环境变量
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config.InitConfig()
}