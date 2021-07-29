package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
	"testing"
	"weather-or-not-bot/internal/service/mock"
)

//go:generate mockgen -source interfaces.go -package=mock -destination=mock/mock.go

func TestMessageService_HandleNewMessage(t *testing.T) {
	ctx := context.Background()

	chatID := int64(123)

	tests := []struct {
		name    string
		prepare func(
			bc *mock.MockBotClient,
			fc *mock.MockForecastClient,
			rf *mock.MockReportFormatter,
			br *mock.MockBotUIRepo,
			lr *mock.MockLocationRepo,
			ulr *mock.MockUserLocationRepo,
			ur *mock.MockUserDataRepo,
			)
		upd *bot.Update
		wantErr bool
	}{
		{
			name:    "1. Error on handling update",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				bc.EXPECT().Send(bot.NewMessage(chatID, "some text")).Return(bot.Message{}, errors.New("some error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bc := mock.NewMockBotClient(ctrl)
			fc := mock.NewMockForecastClient(ctrl)
			rf := mock.NewMockReportFormatter(ctrl)
			ur := mock.NewMockUserDataRepo(ctrl)
			br := mock.NewMockBotUIRepo(ctrl)
			ulr := mock.NewMockUserLocationRepo(ctrl)
			lr := mock.NewMockLocationRepo(ctrl)

			tt.prepare(bc, fc, rf, br, lr, ulr, ur)

			s := NewMessageService(bc, fc, rf, br, lr, ulr, ur)
			if err := s.HandleNewMessage(ctx, tt.upd); (err != nil) != tt.wantErr {
				t.Errorf("HandleNewMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
