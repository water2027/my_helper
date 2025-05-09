package sse

import (
	"fmt"
	"os"
	"slices"
	"time"
	"strconv"

	"wx_assistant/plugins"
	"wx_assistant/plugins/sse/sseapi"
	"wx_assistant/utils"
)

type SsePlugin struct {
	PostChan      chan string
	currentPostId int
	recentPosts   []sseapi.Post
	sseHelper     *sseapi.SSEHelper
}

// 可能存在并发问题，以后哪天需要了再改吧
func (pc *SsePlugin) InitHandler() {
	go func() {
		utils.SetCycleTask(func() {
			posts := pc.sseHelper.GetPosts()
			slices.Reverse(posts)
			for _, post := range posts {
				id := post.PostID
				if id < pc.currentPostId {
					continue
				}
				msg := fmt.Sprintf("%s\nhttps://ssemarket.cn/new/postdetail/%d", post.Title, id)
				pc.PostChan <- msg
				pc.currentPostId = id
				pc.recentPosts = append(pc.recentPosts, post)
			}
		}, time.Minute*1)

	}()
	go func() {
		utils.SetPeriodicTask(func() {
			time.Sleep(time.Second * 5)
			now := time.Now()
			lastTime := now.Add(-time.Hour * 6)
			msg := fmt.Sprintf("%d-%d-%d %d:00:00至%d-%d-%d %d:00:00\n\n", lastTime.Year(), lastTime.Month(), lastTime.Day(), lastTime.Hour(), now.Year(), now.Month(), now.Day(), now.Hour())
			for _, post := range pc.recentPosts {
				msg += fmt.Sprintf("\n\n%s\nhttps://ssemarket.cn/new/postdetail/%d\n\n", post.Title, post.PostID)
			}
			pc.PostChan <- msg
			pc.recentPosts = []sseapi.Post{}
		}, 2025, 5, 9, 7, 0, time.Hour*6)
	}()
}

func (pc *SsePlugin) Name() string {
	return "SSE"
}

func (pc *SsePlugin) GetChan() chan string {
	return pc.PostChan
}

func init() {
	env_id := os.Getenv("SSE_ID")

	// env_id转换为int
	id, err := strconv.Atoi(env_id)
	if err != nil {
		fmt.Println("SSE_ID环境变量转换为int失败:", err)
		return
	}

	sp := &SsePlugin{
		PostChan: make(chan string),
		sseHelper: sseapi.NewSSEHelper(),
		currentPostId: id,
		recentPosts:   []sseapi.Post{},
	}
	plugins.RegisterPlugin(sp)
}
