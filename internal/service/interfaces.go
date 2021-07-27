package service

import (
	"context"
	bot "gopkg.in/telegram-bot-api.v4"
	"weather-or-not-bot/internal/types"
)

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
	GetUserRecentLocation(ctx context.Context, userID int) (*types.UserCoordinates, error)
	SaveLocationName(ctx context.Context, userID int, locationName string) error
	GetCoordinatesByCityName(ctx context.Context, locationName string) (bot.Location, error)
}

type ForecastClient interface {
	GetForecast(ctx context.Context, loc *types.UserCoordinates, period string) ([]byte, error)
}

type ReportFormatter interface {
	FormatNow(ctx context.Context) string
	FormatHours(ctx context.Context, hours int) string
	FormatDays(ctx context.Context, days int) string
}
