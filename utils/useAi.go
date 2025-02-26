package utils

import (
	"fmt"
	"log"
	"os"
)

type AiResponse struct {
	Created   int64    `json:"created"`
	ID        string   `json:"id"`
	Model     string   `json:"model"`
	RequestID string   `json:"request_id"`
	Choices   []Choice `json:"choices"`
	Usage     Usage    `json:"usage"`
}

// Choice represents an individual choice in the response.
type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	Message      Message `json:"message"`
}

// Message represents the message details within a choice.
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// Usage represents the token usage details.
type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type AiBot struct {
	BaseUrl string
	Prompt  string
}

func NewAiBot(url string, prompt string) *AiBot {
	return &AiBot{
		BaseUrl: url,
		Prompt:  prompt,
	}
}

func (ai *AiBot) Send(msg string) string {
	var resp AiResponse
	reqHandler := UseRequest()
	body := fmt.Sprintf(`{"model":"glm-4-flash", "messages":[{"role":"system","content":"%s"},{"role":"user", "content":"%s"}]}`, ai.Prompt, msg)
	api_key := os.Getenv("AI_KEY")
	err := reqHandler.Post(ai.BaseUrl, RequestInit{Body: body, Header: map[string]string{"Authorization":api_key}}, &resp)
	if err != nil {
		log.Println(err)
		return ""
	}
	return resp.Choices[0].Message.Content
}