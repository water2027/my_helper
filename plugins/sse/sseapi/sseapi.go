package sseapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type loginResponse struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Post struct {
	PostID        int    `json:"PostID"`
	UserID        int    `json:"UserID"`
	UserName      string `json:"UserName"`
	UserScore     int    `json:"UserScore"`
	UserTelephone string `json:"UserTelephone"`
	UserAvatar    string `json:"UserAvatar"`
	UserIdentity  string `json:"UserIdentity"`
	Title         string `json:"Title"`
	Content       string `json:"Content"`
	Like          int    `json:"Like"`
	Comment       int    `json:"Comment"`
	Browse        int    `json:"Browse"`
	Heat          int    `json:"Heat"`
	PostTime      string `json:"PostTime"`
	IsSaved       bool   `json:"IsSaved"`
	IsLiked       bool   `json:"IsLiked"`
	Photos        string `json:"Photos"`
	Tag           string `json:"Tag"`
}

//type RatingPost struct {
//	PostID        int    `json:"PostID"`
//	UserID        int    `json:"UserID"`
//	UserName      string `json:"UserName"`
//	UserScore     int    `json:"UserScore"`
//	UserTelephone string `json:"UserTelephone"`
//	UserAvatar    string `json:"UserAvatar"`
//	UserIdentity  string `json:"UserIdentity"`
//	Title         string `json:"Title"`
//	Content       string `json:"Content"`
//	Like          int    `json:"Like"`
//	Comment       int    `json:"Comment"`
//	Browse        int    `json:"Browse"`
//	Heat          int    `json:"Heat"`
//	PostTime      string `json:"PostTime"`
//	IsSaved       bool   `json:"IsSaved"`
//	IsLiked       bool   `json:"IsLiked"`
//	Photos        string `json:"Photos"`
//	Tag           string `json:"Tag"`
//}

func GetPosts() []Post {
	client := &http.Client{}
	loginReq, err := loginSSEReq()
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	req, err := getPostsReq()
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	ratingreq, err := getRatingPostsReq()

	if err != nil {
		log.Println(err)
		return []Post{}
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return []Post{}
	}

	var loginResponse loginResponse

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
	ratingreq.Header.Add("Authorization", "Bearer "+loginResponse.Data.Token)

	defer loginResp.Body.Close()

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	defer resp.Body.Close()

	ratingresp, err := client.Do(ratingreq)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	defer ratingresp.Body.Close()

	var posts []Post
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	json.Unmarshal(body, &posts)

	body, err = io.ReadAll(ratingresp.Body)
	if err != nil {
		log.Println(err)
		return []Post{}
	}
	var ratingposts []Post
	json.Unmarshal(body, &ratingposts)

	posts = append(posts, ratingposts...)
	return posts
}
