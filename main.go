package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func main() {
	//parsing config from .env
	var cfg = parseConfig()

	//establishing connection to database
	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
	if err != nil {
		err = errors.Wrap(err, "Unable to connect to database")
		fmt.Println(err)
	}
	defer conn.Close(context.Background())

	//creating a new BotAPI instance using token
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		err = errors.Wrap(err, "failed to create new BotAPI with given token")
		fmt.Println(err)
	}
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	// bot.Debug = true

	// setting a webhook
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(cfg.WebhookURL))
	if err != nil {
		err = errors.Wrap(err, "failed to set WebHook")
		fmt.Println(err)
	}

	// registering an http handler for a webhook, returns UpdatesChannel (<-chan Update)
	updates := bot.ListenForWebhook("/")

	// launching a server on a local host :8080
	go http.ListenAndServe(cfg.Port, nil)
	fmt.Println("start listen", cfg.Port)

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
