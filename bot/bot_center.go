package bot

import (
	"context"
	"log"

	"wx_assistant/config"
	"wx_assistant/message"
	"wx_assistant/plugins"
)

type botCenter struct {
	bots         map[string](*bot)
	messageChan  chan message.Message
	infoHandlers []plugins.PluginHandlerOption
	cancelFunc   context.CancelFunc
}

func NewBotCenter(bots []config.BotConfig, infoHandlers []plugins.PluginHandlerOption) *botCenter {
	botsMap := make(map[string]*bot)
	for _, botConfig := range bots {
		_bot := newBot(botConfig.Name, botConfig.Webhook)
		botsMap[botConfig.Name] = _bot
	}
	return &botCenter{
		bots:         botsMap,
		infoHandlers: infoHandlers,
		messageChan:  make(chan message.Message),
	}
}

func (bc *botCenter) receiveMessage(ctx context.Context) {
	for _, handler := range bc.infoHandlers {
		h := handler.(plugins.Plugin)
		go func(ctx context.Context) {
			c := h.(plugins.PluginHandlerOption).GetChan()
			for {
				select {
				case msg := <-c:
					{
						bc.messageChan <- msg
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

func (bc *botCenter) handleMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-bc.messageChan:
			from := msg.GetFrom()
			target := msg.GetTarget()
			content := msg.GetContent()
			bot, exists := bc.bots[target]
			if !exists {
				continue
			}
			err := bot.SendMessage(content)
			if err != nil {
				log.Printf("from plugin %s, to %s, content %s, err %s\n", from, target, content, err)
			}
		}
	}
}

func (bc *botCenter) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	bc.cancelFunc = cancel
	defer cancel()
	bc.receiveMessage(ctx)
	return bc.handleMessage(ctx)
}

func (bc *botCenter) Stop() {
	bc.cancelFunc()
}
