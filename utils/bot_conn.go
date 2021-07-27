package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	bot "gopkg.in/telegram-bot-api.v4"
)

func NewBotApi() *bot.BotAPI {
	botAPI, err := bot.NewBotAPI(viper.GetString("bot_token"))
	if err != nil {
		logrus.WithError(err).Fatal("failed to create new BotAPI with given token")
	}

	botAPI.Debug = viper.GetBool("bot_debug_on")
	logrus.Infof("Authorized on account '%s'", botAPI.Self.UserName)

	// setting a webhook
	_, err = botAPI.SetWebhook(bot.NewWebhook(viper.GetString("webhook")))
	if err != nil {
		logrus.WithError(err).Fatal("failed to set WebHook")
	}

	return botAPI
}
