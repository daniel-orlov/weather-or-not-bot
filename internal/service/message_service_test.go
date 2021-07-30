package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
	"testing"
	"weather-or-not-bot/internal/service/mock"
	"weather-or-not-bot/internal/types"
)

//go:generate mockgen -source interfaces.go -package=mock -destination=mock/mock.go

func TestMessageService_HandleNewMessage(t *testing.T) {
	ctx := context.Background()
	someErr := errors.New("some error")

	chatID := int64(123)
	user := &bot.User{
		ID:           122334,
		FirstName:    "John",
		LastName:     "Doe",
		UserName:     "the_john",
		LanguageCode: "en",
		IsBot:        false,
	}

	uLoc := &types.UserCoordinates{
		LocationID: 31415,
		Latitude:   "12,32",
		Longitude:  "45,16",
	}

	// Keyboards
	simpleButton := bot.NewKeyboardButton("some_text")
	locationButton := bot.NewKeyboardButtonLocation("some_text")
	rowWithTwoSimpleButtons := bot.NewKeyboardButtonRow(simpleButton, simpleButton)
	rowWithThreeSimpleButtons := bot.NewKeyboardButtonRow(simpleButton, simpleButton, simpleButton)

	mainMenu := bot.NewReplyKeyboard(bot.NewKeyboardButtonRow(locationButton), bot.NewKeyboardButtonRow(simpleButton))
	chPeriod := bot.NewReplyKeyboard(rowWithTwoSimpleButtons, rowWithTwoSimpleButtons)
	chHours := bot.NewReplyKeyboard(rowWithThreeSimpleButtons, rowWithThreeSimpleButtons)
	chDays := bot.NewReplyKeyboard(rowWithThreeSimpleButtons, rowWithThreeSimpleButtons)

	// Updates.
	stop := &bot.Update{UpdateID: 1, Message: &bot.Message{MessageID: 111, Text: Stop, From: user, Chat: &bot.Chat{ID: chatID}}}
	strt := &bot.Update{UpdateID: 2, Message: &bot.Message{MessageID: 112, Text: Start, From: user, Chat: &bot.Chat{ID: chatID}}}
	b2mm := &bot.Update{UpdateID: 3, Message: &bot.Message{MessageID: 113, Text: BackToMainMenu, From: user, Chat: &bot.Chat{ID: chatID}}}
	back := &bot.Update{UpdateID: 4, Message: &bot.Message{MessageID: 114, Text: Back, From: user, Chat: &bot.Chat{ID: chatID}}}
	byHs := &bot.Update{UpdateID: 5, Message: &bot.Message{MessageID: 115, Text: ByHours, From: user, Chat: &bot.Chat{ID: chatID}}}
	byDs := &bot.Update{UpdateID: 6, Message: &bot.Message{MessageID: 116, Text: ByDays, From: user, Chat: &bot.Chat{ID: chatID}}}
	curr := &bot.Update{UpdateID: 7, Message: &bot.Message{MessageID: 117, Text: CurrentWeather, From: user, Chat: &bot.Chat{ID: chatID}}}

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
			name:    "1. Error on handling Stop",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["End"])
				resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
				bc.EXPECT().Send(resp).Return(bot.Message{}, someErr)
			},
			upd: stop,
			wantErr: true,
		},
		{
			name:    "2. Error on adding user to db on Start",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				ur.EXPECT().AddUserIfNotExists(ctx, user).Return(someErr)
			},
			upd: strt,
			wantErr: true,
		},
		{
			name:    "3. Error on getting user's recent location from db on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(nil, someErr)
			},
			upd: curr,
			wantErr: true,
		},
		{
			name:    "4. Error on getting a forecast on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx,uLoc, curr.Message.Text).Return([]byte{}, someErr)
			},
			upd: curr,
			wantErr: true,
		},
		{
			name:    "4. Error on parsing a forecast on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx,uLoc, curr.Message.Text).Return([]byte{}, someErr)
			},
			upd: curr,
			wantErr: true,
		},

		{
			name:    ". Success on handling Stop",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["End"])
				resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: stop,
			wantErr: false,
		},
		{
			name:    ". Success on handling Start",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				ur.EXPECT().AddUserIfNotExists(ctx, user).Return(nil)
				resp := bot.NewMessage(chatID, fmt.Sprintf("%s\n%s", commentsEn["DefaultMessage"], pickASaying(strt.Message.MessageID, sayingsEn)))
				resp.ReplyMarkup = mainMenu
				br.EXPECT().GetMainMenuKeyboard().Return(mainMenu)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: strt,
			wantErr: false,
		},
		{
			name:    ". Success on handling BackToMainMenu",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["ChooseLocation"])
				resp.ReplyMarkup = mainMenu
				br.EXPECT().GetMainMenuKeyboard().Return(mainMenu)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: b2mm,
			wantErr: false,
		},
		{
			name:    ". Success on handling Back",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriodType"])
				resp.ReplyMarkup = chPeriod
				br.EXPECT().GetDaysOrHoursKeyboard().Return(chPeriod)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: back,
			wantErr: false,
		},
		{
			name:    ". Success on handling ByHours",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriod"])
				resp.ReplyMarkup = chHours
				br.EXPECT().GetHoursKeyboard().Return(chHours)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: byHs,
			wantErr: false,
		},
		{
			name:    ". Success on handling ByDays",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo, ) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriod"])
				resp.ReplyMarkup = chDays
				br.EXPECT().GetHoursKeyboard().Return(chDays)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd: byDs,
			wantErr: false,
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
