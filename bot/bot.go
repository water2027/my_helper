package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type InfoHandler interface {
	GetChan() chan string
}

type bot struct {
	webhook      string
	infoHandlers []InfoHandler
	templateStr  string
	messageChan  chan string
	cancelFunc context.CancelFunc
}

type BotHandler interface {
	ReceiveMessage(ctx context.Context)
	SendMessage(resp string) error
	Run() error
	Stop()
}

func NewBot(webhook string, infoHandlers []InfoHandler, templateStr string) *bot {
	return &bot{
		webhook:      webhook,
		infoHandlers: infoHandlers,
		templateStr:  templateStr,
	}
}

func (b *bot) ReceiveMessage(ctx context.Context) {
	for _, handler := range b.infoHandlers {
		go func(ctx context.Context, handler InfoHandler) {
			c := handler.GetChan()
			for {
				select {
				case msg := <-c:
					{
						b.messageChan <- msg
					}
				case <-ctx.Done():
					{
						return
					}
				}
			}
		}(ctx, handler)
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

func (b *bot) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancelFunc = cancel
	defer cancel()
	b.ReceiveMessage(ctx)
	for message := range b.messageChan {
		err := b.SendMessage(message)
		if err != nil {
			cancel()
			return err
		}
	}
	return nil
}

func (b *bot) Stop() {
	b.cancelFunc()
}
