package utils

import (
	"fmt"
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

func SetTodayTask(f func(), hour, minute int) {
	ctx, cancel := context.WithCancel(context.Background())
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	targetTime := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
	duration := time.Until(targetTime)
	time.AfterFunc(duration, func() {
		f()
		cancel()
	})

	<-ctx.Done()
}

func SetPeriodicTask(f func(), startYear int, startMonth time.Month,
	startDay, startHour, startMinute int, interval time.Duration) {
	ctx := context.Background()
	targetTime := time.Date(startYear, startMonth, startDay, startHour, startMinute, 0, 0, time.Local)

	now := time.Now()
	if targetTime.Before(now) {
		elapsed := now.Sub(targetTime)
		intervals := elapsed / interval
		targetTime = targetTime.Add(interval * (intervals + 1))
	}

	go func() {
		for {
			duration := time.Until(targetTime)
			timer := time.NewTimer(duration)
			select {
			case <-timer.C:
				fmt.Println("Executing periodic task at:", targetTime)
				f()
				targetTime = targetTime.Add(interval)
			case <-ctx.Done():
				fmt.Println("Stopping periodic task at:", targetTime)
				timer.Stop()
				return
			}
		}
	}()
}

func SetCycleTask(f func(), interval time.Duration) {
	now := time.Now().Add(1 * time.Minute)
	SetPeriodicTask(f, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), interval)
}
