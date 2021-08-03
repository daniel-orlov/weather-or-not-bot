package main

import (
	"net/http"
	"weather-or-not-bot/internal/repository"
	"weather-or-not-bot/internal/service"
	"weather-or-not-bot/internal/transport"
	"weather-or-not-bot/internal/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	pflag.String("bot_token", `fake_token`, "Token to access Telegram Bot API")
	pflag.String("port", ":8080", "Port to listen to")
	pflag.Bool("bot_debug_on", true, "Turn on bot debug")

	pflag.String("webhook", "https://some-numbers.ngrok.io", "Webhook URL to get updates from bot")
	pflag.String("weather_api_key", `fake_key`, "Client's key to access weather API")

	pflag.String("language", "", "Service language")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	// Creating logger
	ctx, cancelFunc := utils.NewLogger()
	defer cancelFunc()

	// Establishing connection to database.
	db := utils.NewDBFromEnv()
	defer db.Close()

	// Running init SQL
	err := utils.RunInitMigration(ctx, db)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot listen and serve")
	}

	// Initiating all repositories.
	usrRepo := repository.NewUserDataRepo(db)
	locRepo := repository.NewLocationRepo(db)
	usrLocRepo := repository.NewUserLocationRepo(db)
	botUIRepo := repository.NewBotUIRepo()

	// Establishing client connections.
	botClient := service.NewBotCmd(utils.NewBotApi())
	forecastClient := repository.NewForecastClient()

	formatter := service.NewFormatter()

	// Instantiating main service.
	svc := service.NewMessageService(botClient, forecastClient, formatter, botUIRepo, locRepo, usrLocRepo, usrRepo)

	// Launching a server.
	go func() {
		err := http.ListenAndServe(viper.GetString("port"), nil)
		if err != nil {
			logrus.WithError(err).Fatal("Cannot listen and serve")
		}
	}()
	logrus.Infof("Start listen on port %s", viper.GetString("port"))

	//Handling messages from user.
	updatesHandler := transport.NewUpdatesHandler(svc, botClient.ListenForWebhook("/"))
	updatesHandler.HandleUpdates(ctx)
}
