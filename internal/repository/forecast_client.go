package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	bot "gopkg.in/telegram-bot-api.v4"
	"io/ioutil"
	"net/http"
)

type ForecastClient struct {
}

func (c *ForecastClient) GetForecast(ctc context.Context, loc *bot.Location, period string) ([]byte, error) {
	log := ctxlogrus.Extract(ctx)
	log.Debug("Getting forecast data from a third-party provider")

	cfg := parseConfig()
	lat := fmt.Sprint(loc.Latitude)
	long := fmt.Sprint(loc.Longitude)
	baseWeatherURL := "https://weatherbit-v1-mashape.p.rapidapi.com/"
	url := fmt.Sprintf(
		"%v%vlang=%v&lat=%v&lon=%v", baseWeatherURL, forecasts[period], cfg.Language, lat, long,
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
