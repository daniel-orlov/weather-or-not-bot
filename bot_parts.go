package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4"
	es "github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"strconv"
)

type UserMessage struct {
	bot        *tgbotapi.BotAPI
	connection *pgx.Conn
	update     *tgbotapi.Update
}

func NewUserMessage(bot *tgbotapi.BotAPI, conn *pgx.Conn, update *tgbotapi.Update) UserMessage {
	userMsg := UserMessage{
		bot:        bot,
		connection: conn,
		update:     update,
	}
	return userMsg
}

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

var locationKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonLocation(commentsEn["AtMyLocation"]),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["AtADiffPlace"]),
	),
)

var backToMainMenuKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["Back0"]),
	),
)

var daysOrHoursKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["ByHours"]),
		tgbotapi.NewKeyboardButton(commentsEn["ByDays"]),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["Now"]),
		tgbotapi.NewKeyboardButton(commentsEn["Back0"]),
	),
)

var daysKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["3Days"]),
		tgbotapi.NewKeyboardButton(commentsEn["5Days"]),
		tgbotapi.NewKeyboardButton(commentsEn["7Days"]),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["10Days"]),
		tgbotapi.NewKeyboardButton(commentsEn["16Days"]),
		tgbotapi.NewKeyboardButton(commentsEn["Back1"]),
	),
)

var hoursKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["24Hours"]),
		tgbotapi.NewKeyboardButton(commentsEn["48Hours"]),
		tgbotapi.NewKeyboardButton(commentsEn["72Hours"]),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(commentsEn["96Hours"]),
		tgbotapi.NewKeyboardButton(commentsEn["120Hours"]),
		tgbotapi.NewKeyboardButton(commentsEn["Back1"]),
	),
)

var keyboards = map[string]tgbotapi.ReplyKeyboardMarkup{
	"main":   locationKeyboard,
	"period": daysOrHoursKeyboard,
	"days":   daysKeyboard,
	"hours":  hoursKeyboard,
	"back":   backToMainMenuKeyboard,
}

//HANDLERS
func handleStop(um *UserMessage) {
	fmt.Println("EXECUTING: handleStop")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["End"])
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleStart(um *UserMessage) {
	fmt.Println("EXECUTING: handleStart")
	addUserIfNotExists(um.connection, um.update)
	greeting := commentsEn["DefaultMessage"] + "\n" + pickASaying(sayingsEn)
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, greeting)
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleBackToMainMenu(um *UserMessage) {
	fmt.Println("EXECUTING: handleBackToMainMenu")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["ChooseLocation"])
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleBack(um *UserMessage) {
	fmt.Println("EXECUTING: handleBack")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriodType"])
	msg.ReplyMarkup = keyboards["period"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleByDays(um *UserMessage) {
	fmt.Println("EXECUTING: handleByDays")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriod"])
	msg.ReplyMarkup = keyboards["days"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleByHours(um *UserMessage) {
	fmt.Println("EXECUTING: handleByHours")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["ChoosePeriod"])
	msg.ReplyMarkup = keyboards["hours"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
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
		err = es.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = es.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.Data[0].CityName, um.connection, userId)
	repr := wr.formatNow()

	msg := tgbotapi.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["period"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
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
		err = es.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = es.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.CityName, um.connection, userId)
	repr := wr.formatHours(timePeriod)

	msg := tgbotapi.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["hours"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
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
		err = es.Wrap(err, "failed to getForecast")
		fmt.Println(err)
	}

	wr, err := parseWeather(forecast)
	if err != nil {
		err = es.Wrap(err, "failed to parseWeather")
		fmt.Println(err)
	}

	nameMostRecentLocation(wr.CityName, um.connection, userId)
	repr := wr.formatDays(timePeriod)

	msg := tgbotapi.NewMessage(chatId, repr)
	msg.ReplyMarkup = keyboards["days"]
	_, err = um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleLocationByCoords(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByCoords")
	addLocationByCoords(um.connection, um.update.Message)
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
	msg.ReplyMarkup = keyboards["period"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleWeatherElsewhere(um *UserMessage) {
	fmt.Println("EXECUTING: handleWeatherElsewhere")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["DiffPlaceAccepted"])
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardHide{HideKeyboard: true}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleLocationByText(um *UserMessage) {
	fmt.Println("EXECUTING: handleLocationByText")
	loc := retrieveCoordinates(um.connection, um.update.Message.Text)
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["TryAgain"])
	msg.ReplyMarkup = keyboards["back"]
	if loc.Longitude != 0 && loc.Latitude != 0 { //this one makes it impossible to use bot from one place in Ghana
		um.update.Message.Location = &loc
		addLocationByCoords(um.connection, um.update.Message)
		msg = tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["CoordsAccepted"])
		msg.ReplyMarkup = keyboards["period"]
	}
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
		fmt.Println(err)
	}
}
func handleUnknown(um *UserMessage) {
	fmt.Println("EXECUTING: handleUnknown")
	msg := tgbotapi.NewMessage(um.update.Message.Chat.ID, commentsEn["Unknown"])
	msg.ReplyMarkup = keyboards["main"]
	_, err := um.bot.Send(msg)
	if err != nil {
		err = es.Wrap(err, "Unable to Send message")
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
func addUserIfNotExists(conn *pgx.Conn, update *tgbotapi.Update) {
	fmt.Println("EXECUTING: addUserIfNotExists")
	username := update.Message.From.UserName
	userId := fmt.Sprintf("%v", update.Message.From.ID)
	language := update.Message.From.LanguageCode
	sqlQuery := fmt.Sprintf(
		`INSERT INTO users (user_id, username, language)
				VALUES (%v, '%v', '%v')
				ON CONFLICT (user_id) DO NOTHING
				`,
		userId, username, language,
	)
	err := conn.QueryRow(context.Background(), sqlQuery).Scan()
	if err.Error() == "no rows in result set" { //I have doubts about this one. Not sure if this is right
		fmt.Println("SUCCEEDED: QueryRow when adding UserIfNotExists")
	} else if err != nil {
		err = es.Wrap(err, "FAILED: QueryRow when adding UserIfNotExists")
		fmt.Println(err)
	}
}

func addLocationByCoords(conn *pgx.Conn, msg *tgbotapi.Message) {
	fmt.Println("EXECUTING: addLocationByCoords")
	userId := msg.From.ID
	lat := fmt.Sprintf("%v", msg.Location.Latitude)
	long := fmt.Sprintf("%v", msg.Location.Longitude)

	sqlQuery := fmt.Sprintf(
		`INSERT INTO locations ("user", "latitude", "longitude")
				VALUES (%v, %v, %v)
				`,
		userId, lat, long,
	)
	err := conn.QueryRow(context.Background(), sqlQuery).Scan()
	if err.Error() == "no rows in result set" { //I have doubts about this one. Not sure if this is right
		fmt.Println("SUCCEEDED: QueryRow when adding LocationByCoords")
	} else if err != nil {
		err = es.Wrap(err, "FAILED: QueryRow when adding LocationByCoords")
		fmt.Println(err)
	}
}

func nameMostRecentLocation(name string, conn *pgx.Conn, userId int) {

	fmt.Println("EXECUTING: nameMostRecentLocation")
	loc, id := getMostRecentLocation(conn, userId)
	sqlQuery := fmt.Sprintf(
		`INSERT INTO locations (id, "user", latitude, longitude, location)
				VALUES (%v, %v, '%v', '%v', '%v')
				ON CONFLICT (id) DO UPDATE
				SET location = '%v'
				`,
		id, userId, loc.Latitude, loc.Longitude, name, name,
	)
	fmt.Println(sqlQuery)
	err := conn.QueryRow(context.Background(), sqlQuery).Scan()
	if err.Error() == "no rows in result set" { //I have doubts about this one. Not sure if this is right
		fmt.Println("SUCCEEDED: QueryRow when naming MostRecentLocation")
	} else if err != nil {
		err = es.Wrap(err, "FAILED: QueryRow when naming MostRecentLocation")
		fmt.Println(err)
	}
}
func getMostRecentLocation(conn *pgx.Conn, userId int) (tgbotapi.Location, int) {
	/* TESTED, IT WORKS AS EXPECTED
	 */
	fmt.Println("EXECUTING: getMostRecentLocation")
	loc := tgbotapi.Location{}
	sqlQuery := fmt.Sprintf(
		`SELECT id, latitude, longitude
				FROM locations
				WHERE "user" = '%v'
				ORDER BY id DESC
				LIMIT 1
				`,
		userId,
	)
	fmt.Println(sqlQuery)
	var lat, long string
	var id int
	err := conn.QueryRow(context.Background(), sqlQuery).Scan(&id, &lat, &long)
	if err != nil {
		err = es.Wrap(err, "FAILED: QueryRow when adding LocationByCoords")
		fmt.Println(err)
	}
	fmt.Println("LONG", long, "LAT", lat)
	loc.Latitude, err = strconv.ParseFloat(lat, 64)
	if err != nil {
		err = es.Wrap(err, "Latitude parsing failed")
		fmt.Println(err)
	}
	loc.Longitude, err = strconv.ParseFloat(long, 64)
	if err != nil {
		err = es.Wrap(err, "Longitude parsing failed")
		fmt.Println(err)
	}
	return loc, id
}

func retrieveCoordinates(conn *pgx.Conn, location string) tgbotapi.Location {
	fmt.Println("EXECUTING: retrieveCoordinates")
	loc := tgbotapi.Location{}
	sqlQuery := fmt.Sprintf(
		`SELECT latitude, longitude 
				FROM places
				WHERE city = '%v'
				`,
		location,
	)
	fmt.Println(sqlQuery)
	var lat, long string
	err := conn.QueryRow(context.Background(), sqlQuery).Scan(&lat, &long)
	if err != nil {
		err = es.Wrap(err, "FAILED: QueryRow when retrieving Coordinates")
		fmt.Println(err)
	}
	fmt.Println("LONG", long, "LAT", lat)
	loc.Latitude, err = strconv.ParseFloat(lat, 64)
	if err != nil {
		err = es.Wrap(err, "Latitude parsing failed")
		fmt.Println(err)
	}
	loc.Longitude, err = strconv.ParseFloat(long, 64)
	if err != nil {
		err = es.Wrap(err, "Longitude parsing failed")
		fmt.Println(err)
	}
	return loc
}

func fetchEmojis(backup map[string]int) map[string]int {
	fmt.Println("EXECUTING: fetchEmojis")
	cfg := parseConfig()
	//establishing connection to database
	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
	if err != nil {
		err = es.Wrap(err, "Unable to connect to database")
		fmt.Println(err)
	}
	defer conn.Close(context.Background())

	var emojis = make(map[string]int)
	sqlQuery := `SELECT name, code FROM emojis`
	fmt.Println(sqlQuery)
	rows, err := conn.Query(context.Background(), sqlQuery)
	if err != nil {
		err = es.Wrap(err, "FAILED: Query when fetching Emojis")
		fmt.Println(err)
	}

	var name string
	var code int
	for rows.Next() {
		err = rows.Scan(&name, &code)
		if err != nil {
			err = es.Wrap(err, "FAILED: Scanning a Row while fetching Emojis")
			fmt.Println(err)
		}
		emojis[name] = code
	}
	err = rows.Err()
	if err != nil {
		err = es.Wrap(err, "FAILED: Scan/Next a Row while fetching Emojis")
		fmt.Println(err)
	}
	if len(emojis) == 0 {
		return backup
	}
	return emojis
}
