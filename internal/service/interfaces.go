package service

import (
	"context"
	bot "gopkg.in/telegram-bot-api.v4"
	"weather-or-not-bot/internal/types"
)

type BotClient interface {
	Send(msg bot.MessageConfig) (bot.Message, error)
	ListenForWebhook(webhook string) bot.UpdatesChannel
}

type ForecastClient interface {
	GetForecast(ctx context.Context, loc *types.UserCoordinates, period string) (*types.FullWeatherReport, error)
}

type UserDataRepo interface {
	AddUserIfNotExists(ctx context.Context, user *bot.User) error
}

type BotUIRepo interface {
	GetMainMenuKeyboard() bot.ReplyKeyboardMarkup
	GetBackToMainMenuKeyboard() bot.ReplyKeyboardMarkup
	GetDaysOrHoursKeyboard() bot.ReplyKeyboardMarkup
	GetDaysKeyboard() bot.ReplyKeyboardMarkup
	GetHoursKeyboard() bot.ReplyKeyboardMarkup
}

type LocationRepo interface {
	GetCoordinatesByCityName(ctx context.Context, locationName string) (*bot.Location, error)
}

type UserLocationRepo interface {
	GetUserRecentLocation(ctx context.Context, userID int) (*types.UserCoordinates, error)
	SaveUserLocationName(ctx context.Context, userID int, locationName string) error
	AddUserLocationByCoordinates(ctx context.Context, userID int, loc *bot.Location) error
}

type ReportFormatter interface {
	FormatNow(ctx context.Context, report *types.FullWeatherReport) string
	FormatHours(ctx context.Context, report *types.FullWeatherReport, hours int) string
	FormatDays(ctx context.Context, report *types.FullWeatherReport, days int) string
}
