package transport

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"

	bot "gopkg.in/telegram-bot-api.v4"
)

type MessageService interface {
	HandleNewMessage(ctx context.Context, msg *bot.Update) error
}

type UpdatesHandler struct {
	svc MessageService
	upd bot.UpdatesChannel
}

func NewUpdatesHandler(svc MessageService, upd bot.UpdatesChannel) *UpdatesHandler {
	return &UpdatesHandler{svc: svc, upd: upd}
}

// HandleUpdates reads updates from bot's UpdateChannel and handles them.
func (h *UpdatesHandler) HandleUpdates(ctx context.Context) {
	log := ctxlogrus.Extract(ctx)
	log.Info("Starting handling messages from users")

	for update := range h.upd {
		err := h.svc.HandleNewMessage(ctx, &update)
		if err != nil {
			log.WithError(err).Warn("cannot handle user message")
		}
	}
}
