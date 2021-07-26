package transport

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
)

type MessageService interface {
	HandleSomething()
}

type BotServer struct {
	svc MessageService
}

func NewBotServer(svc MessageService) *BotServer {
	return &BotServer{svc: svc}
}

func (s *BotServer) Stop(ctx context.Context, um *UserMessage) {

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
