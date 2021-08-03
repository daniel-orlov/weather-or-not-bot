package repository

import (
	"context"

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

// TODO add places table to db

const getCoordinatesByCityNameQuery = `
	SELECT latitude, longitude 
	FROM places
	WHERE city = $1;
	`

func (r *LocationRepo) GetCoordinatesByCityName(ctx context.Context, locationName string) (*bot.Location, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"location_name": locationName,
	})
	log.Debug("Getting the coordinates of the location")

	loc := bot.Location{}
	err := r.db.GetContext(ctx, &loc, getCoordinatesByCityNameQuery)
	if err != nil {
		return &loc, errors.Wrap(err, "cannot get coordinates by location name")
	}

	return &loc, nil
}
