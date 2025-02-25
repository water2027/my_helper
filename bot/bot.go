package bot

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"wx_assistant/plugins"
)

type bot struct {
	webhook      string
	infoHandlers []plugins.PluginHandlerOption
	messageChan  chan string
	cancelFunc context.CancelFunc
}

type BotHandler interface {
	ReceiveMessage(ctx context.Context)
	SendMessage(resp string) error
	Run() error
	Stop()
}

func NewBot(webhook string, infoHandlers []plugins.PluginHandlerOption) *bot {
	return &bot{
		webhook:      webhook,
		infoHandlers: infoHandlers,
		messageChan: make(chan string),
	}
}

func (b *bot) ReceiveMessage(ctx context.Context) {
	for _, handler := range b.infoHandlers {
		h := handler.(plugins.Plugin)
		go func(ctx context.Context) {
			c := h.(plugins.PluginHandlerOption).GetChan()
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
		}(ctx)
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
	for {
		select {
		case msg, ok := <-b.messageChan:
			if !ok {
				return nil
			}
			if err := b.SendMessage(msg); err != nil {
				cancel()
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (b *bot) Stop() {
	b.cancelFunc()
}


