package post

import (
	"fmt"
	"time"
)

type TimeStamp struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type RecentPosts struct {
	posts    []Post
	lastTs TimeStamp
}

func (tp *RecentPosts) AddPost(post Post) {
	tp.posts = append(tp.posts, post)
}

func (tp *RecentPosts) Clear() {
	tp.posts = []Post{}
}

func (tp *RecentPosts) SetLastTime(ts TimeStamp) {
	tp.lastTs = ts
}

func (tp *RecentPosts) GetFrom() string {
	return "SSE"
}

func (tp *RecentPosts) GetTarget() string {
	return "raw"
}

func (tp *RecentPosts) GetContent() string {
	now := time.Now()
	lastTime := time.Date(now.Year(), now.Month(), now.Day(), tp.lastTs.Hour, tp.lastTs.Minute, 0, 0, time.Local)
	if lastTime.After(now) {
		lastTime = lastTime.Add(-24 * time.Hour)
	}
	msg := fmt.Sprintf("%d-%d-%d %d:00:00至%d-%d-%d %d:00:00\n\n", lastTime.Year(), lastTime.Month(), lastTime.Day(), lastTime.Hour(), now.Year(), now.Month(), now.Day(), now.Hour())
	
	for _, post := range tp.posts {
		msg += "\n\n" + post.GetContent() + "\n\n"
	}
	// 如果在加入通道后就清空会导致信息丢失, 放在这里吧
	// 哪天再改
	tp.Clear()

	return msg
}

func NewRecentPosts() *RecentPosts {
	return &RecentPosts{
		posts: []Post{},
	}
}