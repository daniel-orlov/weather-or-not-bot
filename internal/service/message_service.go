package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"math/rand"
	"weather-or-not-bot/internal/types"
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

var handlersEn = map[string]func(um *UserMessage){
	"/stop":                  handleStop,
	"/start":                 handleStart,
	"":                       handleEmpty,
	"Weather at my location": handleLocationByCoords,
	"Weather elsewhere":      handleWeatherElsewhere,
	"< Back":                 handleBackToMainMenu,
	"<< Back":                handleBack,
	"By Days":                handleByDays,
	"By Hours":               handleByHours,
	"Now":                    handleNow,
	"3 days":                 handleDays,
	"5 days":                 handleDays,
	"7 days":                 handleDays,
	"10 days":                handleDays,
	"16 days":                handleDays,
	"24 hours":               handleHours,
	"48 hours":               handleHours,
	"72 hours":               handleHours,
	"96 hours":               handleHours,
	"120 hours":              handleHours,
}

// parseWeather parses JSON into FullWeatherReport,	which then can be used to retrieve weather information.
func parseWeather(weather []byte) (*types.FullWeatherReport, error) {
	data := types.FullWeatherReport{}
	err := json.Unmarshal(weather, &data)
	if err != nil {
		return &types.FullWeatherReport{}, errors.Wrap(err, "cannot unmarshal")
	}

	return &data, nil
}

// Picks a random saying.
func pickASaying(sayings []string) string {
	return sayings[rand.Intn(len(sayings))]
}
