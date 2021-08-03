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
		Latitude:   "12.32",
		Longitude:  "45.16",
	}
	botLoc := &bot.Location{
		Longitude: 45.16,
		Latitude:  12.32,
	}

	wr := &types.FullWeatherReport{
		CityName: "some_location",
		Data: []*types.Stat{
			{
				PartOfDay:        Day,
				CityName:         "some_location",
				DateTime:         "",
				WindDirection:    "",
				SunriseTime:      "",
				SunsetTime:       "",
				RelativeHumidity: 0,
				WindSpeedMs:      0,
				IndexUV:          0,
				Precipitation:    0,
				PressureMb:       0,
				Temperature:      0,
				FeelsLikeTemp:    0,
				HighTemp:         0,
				LowTemp:          0,
				Weather:          types.Weather{},
				CloudCoverage:    0,
				Snow:             0,
				IndexAirQuality:  0,
			},
		},
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
	var (
		stop    = &bot.Update{UpdateID: 11, Message: &bot.Message{MessageID: 111, Text: Stop, From: user, Chat: &bot.Chat{ID: chatID}}}
		start   = &bot.Update{UpdateID: 24, Message: &bot.Message{MessageID: 112, Text: Start, From: user, Chat: &bot.Chat{ID: chatID}}}
		back2mm = &bot.Update{UpdateID: 35, Message: &bot.Message{MessageID: 113, Text: BackToMainMenu, From: user, Chat: &bot.Chat{ID: chatID}}}
		back    = &bot.Update{UpdateID: 42, Message: &bot.Message{MessageID: 114, Text: Back, From: user, Chat: &bot.Chat{ID: chatID}}}
		byHours = &bot.Update{UpdateID: 53, Message: &bot.Message{MessageID: 115, Text: ByHours, From: user, Chat: &bot.Chat{ID: chatID}}}
		byDays  = &bot.Update{UpdateID: 69, Message: &bot.Message{MessageID: 116, Text: ByDays, From: user, Chat: &bot.Chat{ID: chatID}}}
		current = &bot.Update{UpdateID: 72, Message: &bot.Message{MessageID: 117, Text: CurrentWeather, From: user, Chat: &bot.Chat{ID: chatID}}}
		days5   = &bot.Update{UpdateID: 88, Message: &bot.Message{MessageID: 118, Text: FiveDays, From: user, Chat: &bot.Chat{ID: chatID}}}
		hours96 = &bot.Update{UpdateID: 90, Message: &bot.Message{MessageID: 119, Text: NinetySixHours, From: user, Chat: &bot.Chat{ID: chatID}}}
		here    = &bot.Update{UpdateID: 14, Message: &bot.Message{MessageID: 120, Text: WeatherHere, From: user, Chat: &bot.Chat{ID: chatID}, Location: botLoc}}
		there   = &bot.Update{UpdateID: 21, Message: &bot.Message{MessageID: 120, Text: WeatherElsewhere, From: user, Chat: &bot.Chat{ID: chatID}}}
		tallinn = &bot.Update{UpdateID: 33, Message: &bot.Message{MessageID: 121, Text: "Tallinn", From: user, Chat: &bot.Chat{ID: chatID}}}
	)

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
		upd     *bot.Update
		wantErr bool
	}{
		{
			name: "1. Error on handling Stop",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["End"])
				resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
				bc.EXPECT().Send(resp).Return(bot.Message{}, someErr)
			},
			upd:     stop,
			wantErr: true,
		},
		{
			name: "2. Error on adding user to db on Start",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ur.EXPECT().AddUserIfNotExists(ctx, user).Return(someErr)
			},
			upd:     start,
			wantErr: true,
		},
		{
			name: "3. Error on getting user's recent location from db on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(nil, someErr)
			},
			upd:     current,
			wantErr: true,
		},
		{
			name: "4. Error on getting a forecast on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx, uLoc, current.Message.Text).Return(nil, someErr)
			},
			upd:     current,
			wantErr: true,
		},
		{
			name: "5. No error, but failed saving location name on Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx, uLoc, current.Message.Text).Return(wr, nil)
				rf.EXPECT().FormatNow(ctx, wr).Return("formatted_report")
				resp := bot.NewMessage(chatID, "formatted_report")
				resp.ReplyMarkup = chPeriod
				br.EXPECT().GetDaysOrHoursKeyboard().Return(chPeriod)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
				ulr.EXPECT().SaveUserLocationName(ctx, user.ID, wr.CityName).Return(someErr)
			},
			upd:     current,
			wantErr: false,
		},
		{
			name: "6. Error on adding user's location from db on WeatherHere",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().AddUserLocationByCoordinates(ctx, user.ID, here.Message.Location).Return(someErr)
			},
			upd:     here,
			wantErr: true,
		},
		{
			name: "7. No error, but failed to get coordinates from db when handling location by text",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				lr.EXPECT().GetCoordinatesByCityName(ctx, tallinn.Message.Text).Return(&bot.Location{}, someErr)
				resp := bot.NewMessage(chatID, commentsEn["Unknown"])
				resp.ReplyMarkup = mainMenu
				br.EXPECT().GetMainMenuKeyboard().Return(mainMenu)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     tallinn,
			wantErr: false,
		},
		{
			name: "8. Success on handling Stop",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["End"])
				resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     stop,
			wantErr: false,
		},
		{
			name: "9. Success on handling Start",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ur.EXPECT().AddUserIfNotExists(ctx, user).Return(nil)
				resp := bot.NewMessage(chatID, fmt.Sprintf("%s\n%s", commentsEn["DefaultMessage"], pickASaying(start.Message.MessageID, sayingsEn)))
				resp.ReplyMarkup = mainMenu
				br.EXPECT().GetMainMenuKeyboard().Return(mainMenu)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     start,
			wantErr: false,
		},
		{
			name: "10. Success on handling BackToMainMenu",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["ChooseLocation"])
				resp.ReplyMarkup = mainMenu
				br.EXPECT().GetMainMenuKeyboard().Return(mainMenu)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     back2mm,
			wantErr: false,
		},
		{
			name: "11. Success on handling Back",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriodType"])
				resp.ReplyMarkup = chPeriod
				br.EXPECT().GetDaysOrHoursKeyboard().Return(chPeriod)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     back,
			wantErr: false,
		},
		{
			name: "12. Success on handling ByHours",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriod"])
				resp.ReplyMarkup = chHours
				br.EXPECT().GetHoursKeyboard().Return(chHours)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     byHours,
			wantErr: false,
		},
		{
			name: "13. Success on handling ByDays",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["ChoosePeriod"])
				resp.ReplyMarkup = chDays
				br.EXPECT().GetDaysKeyboard().Return(chDays)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     byDays,
			wantErr: false,
		},
		{
			name: "14. Success on handling Now",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx, uLoc, current.Message.Text).Return(wr, nil)
				rf.EXPECT().FormatNow(ctx, wr).Return("formatted_report")
				resp := bot.NewMessage(chatID, "formatted_report")
				resp.ReplyMarkup = chPeriod
				br.EXPECT().GetDaysOrHoursKeyboard().Return(chPeriod)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
				ulr.EXPECT().SaveUserLocationName(ctx, user.ID, wr.CityName).Return(nil)
			},
			upd:     current,
			wantErr: false,
		},
		{
			name: "15. Success on handling FiveDays",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx, uLoc, days5.Message.Text).Return(wr, nil)
				rf.EXPECT().FormatDays(ctx, wr, extractNumerals(days5.Message.Text)).Return("formatted_report")
				resp := bot.NewMessage(chatID, "formatted_report")
				resp.ReplyMarkup = chDays
				br.EXPECT().GetDaysKeyboard().Return(chDays)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
				ulr.EXPECT().SaveUserLocationName(ctx, user.ID, wr.CityName).Return(nil)
			},
			upd:     days5,
			wantErr: false,
		},
		{
			name: "16. Success on handling NinetySixHours",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().GetUserRecentLocation(ctx, user.ID).Return(uLoc, nil)
				fc.EXPECT().GetForecast(ctx, uLoc, hours96.Message.Text).Return(wr, nil)
				rf.EXPECT().FormatHours(ctx, wr, extractNumerals(hours96.Message.Text)).Return("formatted_report")
				resp := bot.NewMessage(chatID, "formatted_report")
				resp.ReplyMarkup = chHours
				br.EXPECT().GetHoursKeyboard().Return(chHours)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
				ulr.EXPECT().SaveUserLocationName(ctx, user.ID, wr.CityName).Return(nil)
			},
			upd:     hours96,
			wantErr: false,
		},
		{
			name: "17. Success on handling WeatherHere",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				ulr.EXPECT().AddUserLocationByCoordinates(ctx, user.ID, here.Message.Location).Return(nil)
				resp := bot.NewMessage(chatID, commentsEn["CoordsAccepted"])
				resp.ReplyMarkup = chPeriod
				br.EXPECT().GetDaysOrHoursKeyboard().Return(chPeriod)
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     here,
			wantErr: false,
		},
		{
			name: "18. Success on handling WeatherElsewhere",
			prepare: func(bc *mock.MockBotClient, fc *mock.MockForecastClient, rf *mock.MockReportFormatter, br *mock.MockBotUIRepo, lr *mock.MockLocationRepo, ulr *mock.MockUserLocationRepo, ur *mock.MockUserDataRepo) {
				resp := bot.NewMessage(chatID, commentsEn["DiffPlaceAccepted"])
				resp.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
				bc.EXPECT().Send(resp).Return(bot.Message{}, nil)
			},
			upd:     there,
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
