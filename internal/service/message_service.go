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

func handleStop(um *UserMessage) {
	fmt.Println("EXECUTING: handleStop")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["End"])
	msg.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleStart(um *UserMessage) {
	fmt.Println("EXECUTING: handleStart")
	addUserIfNotExists(um.connection, um.update)
	greeting := commentsEn["DefaultMessage"] + "\n" + pickASaying(sayingsEn)
	msg := bot.NewMessage(um.update.Message.Chat.ID, greeting)
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleBackToMainMenu(um *UserMessage) {
	fmt.Println("EXECUTING: handleBackToMainMenu")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["ChooseLocation"])
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleBack(um *UserMessage) {
	fmt.Println("EXECUTING: handleBack")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriodType"])
	msg.ReplyMarkup = keyboards["period"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleByDays(um *UserMessage) {
	fmt.Println("EXECUTING: handleByDays")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriod"])
	msg.ReplyMarkup = keyboards["days"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleByHours(um *UserMessage) {
	fmt.Println("EXECUTING: handleByHours")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriod"])
	msg.ReplyMarkup = keyboards["hours"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleNow(um *UserMessage) {
	fmt.Println("EXECUTING: handleNow")
	text := um.update.Message.Text
	userId := um.update.Message.From.ID
	chatId := um.update.Message.Chat.ID

	loc, _ := getMostRecentLocation(um.connection, userId)
	forecast, err := getForecast(&loc, text)
	if err != nil {
		err = errors.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = errors.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.Data[0].CityName, um.connection, userId)
	repr := wr.formatNow()

	msg := bot.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["period"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleHours(um *UserMessage) {
	fmt.Println("EXECUTING: handleHours")
	text := um.update.Message.Text
	timePeriod := timePeriodsEn[text]
	userId := um.update.Message.From.ID
	chatId := um.update.Message.Chat.ID

	loc, _ := getMostRecentLocation(um.connection, userId)
	forecast, err := getForecast(&loc, text)
	if err != nil {
		err = errors.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = errors.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.CityName, um.connection, userId)
	repr := wr.formatHours(timePeriod)

	msg := bot.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["hours"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleDays(um *UserMessage) {
	fmt.Println("EXECUTING: handleDays")
	text := um.update.Message.Text
	timePeriod := timePeriodsEn[text]
	userId := um.update.Message.From.ID
	chatId := um.update.Message.Chat.ID

	loc, _ := getMostRecentLocation(um.connection, userId)
	forecast, err := getForecast(&loc, text)
	if err != nil {
		err = errors.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = errors.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.CityName, um.connection, userId)
	repr := wr.formatDays(timePeriod)

	msg := bot.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["days"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleLocationByCoords(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByCoords")
	addLocationByCoords(um.connection, um.update.Message)
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
	msg.ReplyMarkup = keyboards["period"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleWeatherElsewhere(um *UserMessage) {
	fmt.Println("EXECUTING: handleWeatherElsewhere")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["DiffPlaceAccepted"])
	msg.ReplyMarkup = bot.ReplyKeyboardHide{HideKeyboard: true}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleLocationByText(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByText")
	loc := retrieveCoordinates(um.connection, um.update.Message.Text)
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["TryAgain"])
	msg.ReplyMarkup = keyboards["back"]
	if loc.Longitude != 0 && loc.Latitude != 0 { //this one makes it impossible to use bot from one place in Ghana
		um.update.Message.Location = &loc
		addLocationByCoords(um.connection, um.update.Message)
		msg = bot.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
		msg.ReplyMarkup = keyboards["period"]
	}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleUnknown(um *UserMessage) {
	fmt.Println("EXECUTING: handleUnknown")
	msg := bot.NewMessage(um.update.Message.Chat.ID, commentsEn["Unknown"])
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = errors.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleEmpty(um *UserMessage) {
	fmt.Println("EXECUTING: handleEmpty")
	if um.update.Message.Location != nil {
		handleLocationByCoords(um)
	} else {
		handleUnknown(um)
	}
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
