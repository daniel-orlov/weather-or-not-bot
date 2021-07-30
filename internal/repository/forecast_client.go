package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

const (
	BaseURL = "https://weatherbit-v1-mashape.p.rapidapi.com/"
	HostHeader = "weatherbit-v1-mashape.p.rapidapi.com"
)

func (c *ForecastClient) GetForecast(ctx context.Context, loc *types.UserCoordinates, period string) (*types.FullWeatherReport, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"period": period,
	})
	log.Debug("Getting forecast data from a third-party provider")

	rawWR, err := c.getForecast(loc, period)


}
func (c *ForecastClient) getForecast(loc *types.UserCoordinates, period string) ([]byte, error) {
	url := fmt.Sprintf(
		"%v%vlang=%v&lat=%v&lon=%v", BaseURL, forecasts[period], viper.GetString("language"), loc.Latitude, loc.Longitude,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a new request")
	}
	req.Header.Add("x-rapidapi-host", HostHeader)
	req.Header.Add("x-rapidapi-key", viper.GetString("weather_api_key"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot perform request")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read from response body")
	}

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
