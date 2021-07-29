package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	bot "gopkg.in/telegram-bot-api.v4"
)

// NewBotApi creates a BotAPI instance from token and sets a webhook.
func NewBotApi() *bot.BotAPI {
	botAPI, err := bot.NewBotAPI(viper.GetString("bot_token"))
	if err != nil {
		logrus.WithError(err).Fatal("Cannot create new BotAPI with given token")
	}

	botAPI.Debug = viper.GetBool("bot_debug_on")
	logrus.Infof("Authorized on account '%s'", botAPI.Self.UserName)

	// Deleting existing webhook from bot
	_, err = botAPI.RemoveWebhook()
	if err != nil {
		logrus.WithError(err).Fatal("Cannot delete WebHook")
	}

	// Setting a webhook (e.g. using ngrok)
	_, err = botAPI.SetWebhook(bot.NewWebhook(viper.GetString("webhook")))
	if err != nil {
		logrus.WithError(err).Fatal("Cannot set webhook")
	}

	return botAPI
}
