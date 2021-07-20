package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
)

//envDefault:"./.env"

type config struct {
	BotToken   string `env:"BOT_TOKEN"`
	Port       string `env:"PORT"`
	WebhookURL string `env:"WEBHOOK"`
	WeatherAPI string `env:"WEATHER_API"`
	Language   string `env:"LANGUAGE"`
	DbUrl      string `env:"DATABASE_URL"`
}

func parseConfig() config {
	fmt.Println("EXECUTING: parseConfig")
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	//fmt.Printf("%+v\n", cfg)
	return cfg
}

var locationKeyboard = bot.NewReplyKeyboard(
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButtonLocation(commentsEn["AtMyLocation"]),
	),
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["AtADiffPlace"]),
	),
)

var backToMainMenuKeyboard = bot.NewReplyKeyboard(
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["Back0"]),
	),
)

var daysOrHoursKeyboard = bot.NewReplyKeyboard(
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["ByHours"]),
		bot.NewKeyboardButton(commentsEn["ByDays"]),
	),
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["Now"]),
		bot.NewKeyboardButton(commentsEn["Back0"]),
	),
)

var daysKeyboard = bot.NewReplyKeyboard(
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["3Days"]),
		bot.NewKeyboardButton(commentsEn["5Days"]),
		bot.NewKeyboardButton(commentsEn["7Days"]),
	),
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["10Days"]),
		bot.NewKeyboardButton(commentsEn["16Days"]),
		bot.NewKeyboardButton(commentsEn["Back1"]),
	),
)

var hoursKeyboard = bot.NewReplyKeyboard(
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["24Hours"]),
		bot.NewKeyboardButton(commentsEn["48Hours"]),
		bot.NewKeyboardButton(commentsEn["72Hours"]),
	),
	bot.NewKeyboardButtonRow(
		bot.NewKeyboardButton(commentsEn["96Hours"]),
		bot.NewKeyboardButton(commentsEn["120Hours"]),
		bot.NewKeyboardButton(commentsEn["Back1"]),
	),
)

var keyboards = map[string]bot.ReplyKeyboardMarkup{
	"main":   locationKeyboard,
	"period": daysOrHoursKeyboard,
	"days":   daysKeyboard,
	"hours":  hoursKeyboard,
	"back":   backToMainMenuKeyboard,
}

//HANDLERS
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

//DATABASE CALLS

func fetchEmojis(backup map[string]int) map[string]int {
	fmt.Println("EXECUTING: fetchEmojis")
	cfg := parseConfig()
	//establishing connection to database
	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to database")
		fmt.Println(err)
	}
	defer conn.Close(context.Background())

	var emojis = make(map[string]int)
	sqlQuery := `SELECT name, code FROM emojis`
	fmt.Println(sqlQuery)
	rows, err := conn.Query(context.Background(), sqlQuery)
	if err != nil {
		err = errors.Wrap(err, "FAILED: Query when fetching Emojis")
		fmt.Println(err)
	}

	var name string
	var code int
	for rows.Next() {
		err = rows.Scan(&name, &code)
		if err != nil {
			err = errors.Wrap(err, "FAILED: Scanning a Row while fetching Emojis")
			fmt.Println(err)
		}
		emojis[name] = code
	}
	err = rows.Err()
	if err != nil {
		err = errors.Wrap(err, "FAILED: Scan/Next a Row while fetching Emojis")
		fmt.Println(err)
	}
	if len(emojis) == 0 {
		return backup
	}
	return emojis
}
