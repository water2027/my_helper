package main

import (
	"log"
	"os"
	"fmt"
	"os/signal"
	"syscall"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"
)

func initLog() {
	logFile, err := os.OpenFile("ssebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	log.SetOutput(logFile)

	// 创建一个信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 来监听信号
	go func() {
		sig := <-sigChan
		log.Printf("接收到信号: %v", sig)
		// 在这里可以添加其他清理操作
		os.Exit(0)
	}()
}

func main() {
	initLog()
	config.InitConfig()
	infoHandlers := plugins.GetHandlers()
	b := bot.NewBot(config.MyConfig.BotConfig.Webhook, infoHandlers)
	err := b.Run()
	if err != nil {
		log.Println(err)
	}
}
