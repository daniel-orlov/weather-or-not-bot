package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type config struct {
	WebhookURL string `env:"WEBHOOK"`
	WeatherAPI string `env:"WEATHER_API"`
	Language   string `env:"LANGUAGE"`
	DbUrl      string `env:"DATABASE_URL"`
}

func init() {
	pflag.String("bot_token", `fake_token`, "Token to access Telegram Bot API")
	pflag.String("port", "8080", "Port to listen to")

	pflag.String("webhook", "", "Webhook URL to get weather forecasts from")
	pflag.String("weather_api_key", "", "Client's key to access weather API")

	pflag.String("database_url", "", "Database access URL")
	pflag.String("language", "", "Service language")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	//parsing config from .env
	var cfg = parseConfig()

	//establishing connection to database
	conn, err := pgx.Connect(context.Background(), viper.GetString("database_url"))
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to database")
		fmt.Println(err)
	}
	defer conn.Close(context.Background())

	//creating a new BotAPI instance using token
	bot, err := tgbotapi.NewBotAPI(viper.GetString("bot_token"))
	if err != nil {
		err = errors.Wrap(err, "failed to create new BotAPI with given token")
		fmt.Println(err)
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	// bot.Debug = true

	// setting a webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(viper.GetString("webhook")))
	if err != nil {
		err = errors.Wrap(err, "failed to set WebHook")
		fmt.Println(err)
	}

	// registering an http handler for a webhook, returns UpdatesChannel (<-chan Update)
	updates := bot.ListenForWebhook("/")

	// launching a server on a local host :8080
	go http.ListenAndServe(viper.GetString("port"), nil)
	fmt.Println("start listen", viper.GetString("port"))

	// handling messages from user
	for update := range updates {
		fmt.Printf("TEXT: %+v\n", update.Message.Text)
		userMsg := NewUserMessage(bot, conn, &update)
		text := update.Message.Text
		if handle, ok := handlersEn[text]; ok {
			handle(&userMsg)
		} else {
			handleLocationByText(&userMsg)
		}
	}
}
