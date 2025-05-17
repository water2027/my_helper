package main

import (
	"log"
	"os"
	"fmt"
	"os/signal"
	"syscall"

	_ "wx_assistant/init"

	"wx_assistant/bot"
	"wx_assistant/config"
	"wx_assistant/plugins"

	_ "wx_assistant/plugins/sse"
)

func initLog() {
	logFile, err := os.OpenFile("ssebot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Failed to open error log file:", err)
		return
	}
	log.SetOutput(logFile)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("接收到信号: %v", sig)
		os.Exit(0)
	}()
}

func main() {
	initLog()
	infoHandlers := plugins.GetHandlers()
	b := bot.NewBotCenter(config.MyConfig.BotConfig, infoHandlers)
	err := b.Run()
	if err != nil {
		log.Println(err)
	}
}
