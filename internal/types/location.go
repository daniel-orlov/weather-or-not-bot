package types

import (
	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
	"strconv"
)

type WorldCity struct {
	Latitude   string `db:"lat"`
	Longitude  string `db:"long"`
}

func (wc *WorldCity) WorldCityToBotLocation() (*bot.Location, error) {
	lat, err := strconv.ParseFloat(wc.Latitude, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse latitude '%s' from db", wc.Latitude)
	}

	long, err := strconv.ParseFloat(wc.Longitude, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse longitude '%s' from db", wc.Longitude)
	}

	return &bot.Location{
		Longitude: long,
		Latitude:  lat,
	}, nil
}