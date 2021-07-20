package service

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type UserDataRepo interface {
}

type MessageService struct {
	botClient *tgbotapi.BotAPI
	userRepo  UserDataRepo
	update    *tgbotapi.Update
}

func NewMessageService(bot *tgbotapi.BotAPI, repo UserDataRepo, update *tgbotapi.Update) *MessageService {
	return &MessageService{
		botClient: bot,
		userRepo:  repo,
		update:    update,
	}
}

func (s MessageService) SendMessage(ctx context.Context) {
	log := ctxlogrus.Extract(ctx)
	log.Trace("sending message")
}
