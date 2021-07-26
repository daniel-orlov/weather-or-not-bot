package types

import (
	"fmt"
	"github.com/caarlos0/env"
)

// TODO change this to init() in main w/ viper
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

	return cfg
}
