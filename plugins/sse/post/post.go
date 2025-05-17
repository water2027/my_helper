package post

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"wx_assistant/plugins/sse/sseapi"
)

type PostGenerator struct {
	Email     string
	Password  string
	Telephone string
}

func NewGenerator() *PostGenerator {
	email := os.Getenv("SSE_EMAIL")
	password := os.Getenv("SSE_PASSWORD")
	telephone := os.Getenv("SSE_TELEPHONE")
	if email == "" || password == "" || telephone == "" {
		panic("SSE credentials not set in environment variables")
	}
	return &PostGenerator{
		Email:     email,
		Password:  password,
		Telephone: telephone,
	}
}

func (h *PostGenerator) GetPosts() []Post {
	client := &http.Client{}
	loginReq, err := sseapi.LoginSSEReq(h.Email, h.Password)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	req, err := sseapi.GetPostsReq(h.Telephone)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	var loginResponse sseapi.LoginResponse

	body, _ := io.ReadAll(loginResp.Body)
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	if loginResponse.Code != 200 {
		log.Println(loginResponse.Msg)
		return []Post{}
	}
	// 将token添加到第二个请求的header中
	req.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)

	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	defer resp.Body.Close()

	var posts []Post
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	json.Unmarshal(body, &posts)
	return posts
}

type PostTarget int

const (
	RawPost PostTarget = iota
	MarkdownPost
)

type Post struct {
	PostID int `json:"PostID"`
	// UserID        int    `json:"UserID"`
	// UserName      string `json:"UserName"`
	// UserScore     int    `json:"UserScore"`
	// UserTelephone string `json:"UserTelephone"`
	// UserAvatar    string `json:"UserAvatar"`
	// UserIdentity  string `json:"UserIdentity"`
	Title string `json:"Title"`
	// Content       string `json:"Content"`
	// Like          int    `json:"Like"`
	// Comment       int    `json:"Comment"`
	// Browse        int    `json:"Browse"`
	// Heat          int    `json:"Heat"`
	// PostTime      string `json:"PostTime"`
	// IsSaved       bool   `json:"IsSaved"`
	// IsLiked       bool   `json:"IsLiked"`
	// Photos        string `json:"Photos"`
	// Tag           string `json:"Tag"`
	Target PostTarget
}

func (p *Post) GetFrom() string {
	return "SSE"
}

func (p *Post) GetTarget() string {
	switch p.Target {
	case RawPost:
		return "raw"
	case MarkdownPost:
		return "markdown"
	default:
		return "unknown"
	}
}

func (p *Post) GetContent() string {
	switch p.Target {
	case RawPost:
		return fmt.Sprintf("%s\nhttps://ssemarket.cn/new/postdetail/%d", p.Title, p.PostID)
	case MarkdownPost:
		return fmt.Sprintf("[%s](https://ssemarket.cn/new/postdetail/%d)", p.Title, p.PostID)
	default:
		return fmt.Sprintf("%s\nhttps://ssemarket.cn/new/postdetail/%d", p.Title, p.PostID)
	}
}


