package sse

import (
	"context"
	"fmt"
	"log"
	"slices"

	"wx_assistant/plugins"
	"wx_assistant/plugins/sse/sseapi"
	"wx_assistant/utils"
)

type SsePlugin struct {
	PostChan chan string
}

type PostHandler interface {
	Get(ctx context.Context) error
	GetChan() *chan sseapi.Post
}

func (pc *SsePlugin) InitHandler() {
	log.Println("SSE plugin init")
	utils.SetOnceTask(func() {
		posts := sseapi.GetPosts()
		slices.Reverse(posts)
		for _, post := range posts {
			msg := fmt.Sprintf("%s, https://ssemarket.cn/postdetail/%d", post.Title, post.PostID)
			pc.PostChan <- msg
		}
	}, 2025, 2, 24, 9, 29)
}

func (pc *SsePlugin) Name() string {
	return "SSE"
}

func (pc *SsePlugin) GetChan() chan string {
	return pc.PostChan
}

func init() {
	sp := &SsePlugin{
		PostChan: make(chan string),
	}
	plugins.RegisterPlugin(sp)
}
