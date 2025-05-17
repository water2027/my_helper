package sse

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"

	"wx_assistant/message"
	"wx_assistant/plugins"
	"wx_assistant/plugins/sse/post"
	"wx_assistant/utils"
)

type SsePlugin struct {
	PostChan      chan message.Message
	currentPostId int
	recentPosts   post.RecentPosts
	sseHelper     *post.PostGenerator
	timeMap       []post.TimeStamp
}

func (pc *SsePlugin) setTask(lastTs, currentTs post.TimeStamp) {
	utils.SetPeriodicTask(func() {
		pc.recentPosts.SetLastTime(lastTs)
		pc.PostChan <- &pc.recentPosts
	}, 2025, 5, 10, currentTs.Hour, currentTs.Minute, time.Hour*24)
}

// 可能存在并发问题，以后哪天需要了再改吧
func (pc *SsePlugin) InitHandler() {
	go func() {
		utils.SetCycleTask(func() {
			posts := pc.sseHelper.GetPosts()
			slices.Reverse(posts)
			for _, postItem := range posts {
				id := postItem.PostID
				if id <= pc.currentPostId {
					continue
				}
				pc.PostChan <- &postItem
				pc.currentPostId = id
				pc.recentPosts.AddPost(postItem)
				p := postItem
				p.Target = post.MarkdownPost
				pc.PostChan <- &p
			}
		}, time.Minute*1)
	}()
	for i := range pc.timeMap {
		if i == 0 {
			go pc.setTask(pc.timeMap[len(pc.timeMap)-1], pc.timeMap[i])
			continue
		}
		go pc.setTask(pc.timeMap[i-1], pc.timeMap[i])
	}

}

func (pc *SsePlugin) Name() string {
	return "SSE"
}

func (pc *SsePlugin) GetChan() chan message.Message {
	return pc.PostChan
}

func init() {
	env_id := os.Getenv("SSE_ID")

	id, err := strconv.Atoi(env_id)
	if err != nil {
		fmt.Println("SSE_ID环境变量转换为int失败:", err)
		return
	}

	sp := &SsePlugin{
		PostChan:      make(chan message.Message),
		sseHelper:     post.NewGenerator(),
		currentPostId: id,
		recentPosts:   *post.NewRecentPosts(),
		timeMap: []post.TimeStamp{
			{Hour: 8, Minute: 0},
			{Hour: 12, Minute: 0},
			{Hour: 18, Minute: 0},
			{Hour: 22, Minute: 00},
		},
	}
	plugins.RegisterPlugin(sp)
}
