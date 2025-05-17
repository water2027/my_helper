package bot

import (
	"bytes"
	"fmt"
	"net/http"
)

type bot struct {
	name    string
	webhook string
}

func newBot(name, webhook string) *bot {
	return &bot{
		name:    name,
		webhook: webhook,
	}
}

func (b *bot) SendMessage(resp string) error {
	data := fmt.Sprintf(`{"msgtype":"markdown","markdown":{"content":"%s"}}`, resp)
	req, err := http.NewRequest("POST", b.webhook, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

