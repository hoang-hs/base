package alert

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/hoang-hs/base/common/log"
	"os"
)

type Sender interface {
	SendMessage(ctx context.Context, text string)
}

const (
	sendMessagePath = "SendMessage"
)

type Telegram struct {
	client          *resty.Client
	chatID          string
	messageThreadID string
}

func NewTelegram() Sender {
	var (
		telegramToken           = os.Getenv("TELEGRAM_TOKEN")
		telegramChatID          = os.Getenv("TELEGRAM_CHAT_ID")
		telegramMessageThreadID = os.Getenv("TELEGRAM_MESSAGE_THREAD_ID")
	)
	if telegramToken == "" || telegramChatID == "" || telegramMessageThreadID == "" {
		log.Warn("Telegram token, chat id, or message thread is empty")
		return NewNoop()
	}

	return &Telegram{
		client:          resty.New().SetBaseURL("https://api.telegram.org/bot" + telegramToken),
		chatID:          telegramChatID,
		messageThreadID: telegramMessageThreadID,
	}
}

func (a *Telegram) SendMessage(ctx context.Context, text string) {
	resp, err := a.client.R().SetContext(ctx).
		SetQueryParams(map[string]string{
			"chat_id":           a.chatID,
			"text":              text,
			"message_thread_id": a.messageThreadID,
		}).
		SetHeader("Accept", "application/json").
		Get(sendMessagePath)
	if err != nil {
		log.Error("failed to send telegram bot", log.Err(err))
		return
	}
	if resp.IsError() {
		log.Error("failed to send telegram bot", log.String("resp", resp.String()))
	}
}

func NewNoop() *Noop {
	return &Noop{}
}

type Noop struct{}

func (a *Noop) SendMessage(_ context.Context, text string) {
	log.Info("Noop telegram bot", log.String("text", text))
}
