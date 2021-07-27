package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"weather-or-not-bot/internal/types"
)

type ForecastClient struct {
}

func NewForecastClient() *ForecastClient {
	return &ForecastClient{}
}

// TODO make this pretty and refactor

func (c *ForecastClient) GetForecast(ctx context.Context, loc *types.UserCoordinates, period string) ([]byte, error) {
	log := ctxlogrus.Extract(ctx)
	log.Debug("Getting forecast data from a third-party provider")

	lat := fmt.Sprint(loc.Latitude)
	long := fmt.Sprint(loc.Longitude)
	baseWeatherURL := "https://weatherbit-v1-mashape.p.rapidapi.com/"
	url := fmt.Sprintf(
		"%v%vlang=%v&lat=%v&lon=%v", baseWeatherURL, forecasts[period], viper.GetString("language"), lat, long,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("x-rapidapi-host", "weatherbit-v1-mashape.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", viper.GetString("weather_api_key"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to perform request")
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read from response body")
		return nil, err
	}
	fmt.Println("BODY:", string(body))
	return body, nil
}

var forecasts = map[string]string{
	"Now":       "current?",
	"3 days":    "forecast/daily?",
	"5 days":    "forecast/daily?",
	"7 days":    "forecast/daily?",
	"10 days":   "forecast/daily?",
	"16 days":   "forecast/daily?",
	"24 hours":  "forecast/hourly?hours=24&",
	"48 hours":  "forecast/hourly?hours=48&",
	"72 hours":  "forecast/hourly?hours=72&",
	"96 hours":  "forecast/hourly?hours=96&",
	"120 hours": "forecast/hourly?hours=120&",
}
