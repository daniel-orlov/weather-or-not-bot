package repository

import (
	"context"
	"weather-or-not-bot/internal/types"

	bot "gopkg.in/telegram-bot-api.v4"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LocationRepo struct {
	db *sqlx.DB
}

func NewLocationRepo(db *sqlx.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

const getCoordinatesByCityNameQuery = `
	SELECT lat, long
	FROM world_cities
	WHERE city_ascii = $1;
	`

func (r *LocationRepo) GetCoordinatesByCityName(ctx context.Context, locationName string) (*bot.Location, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"location_name": locationName,
	})
	log.Debug("Getting the coordinates of the location")

	cityLoc := types.WorldCity{}
	err := r.db.GetContext(ctx, &cityLoc, getCoordinatesByCityNameQuery, locationName)
	if err != nil {
		return &bot.Location{}, errors.Wrap(err, "cannot get coordinates by location name")
	}

	return cityLoc.WorldCityToBotLocation()
}
