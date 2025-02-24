package utils

import (
	"time"

	"golang.org/x/net/context"
)

func SetOnceTask(f func(), year int, month time.Month, day, hour, minute int) {
	ctx, cancel := context.WithCancel(context.Background())
	targetTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
	duration := time.Until(targetTime)
	time.AfterFunc(duration, func() {
		f()
		cancel()
	})

	<-ctx.Done()
}