package repository

import (
	"context"
	bot "gopkg.in/telegram-bot-api.v4"
	"weather-or-not-bot/internal/types"

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

func (r LocationRepo) GetUserRecentLocation(ctx context.Context, userID int) (*types.UserCoordinates, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id": userID,
	})
	log.Debug("getting user's recent coordinates from db")

	userLocation := types.UserCoordinates{}
	err := r.db.GetContext(ctx, &userLocation, getUserRecentLocationQuery)
	if err != nil {
		return &userLocation, errors.Wrap(err, "cannot get user's recent location")
	}

	return &userLocation, nil
}

const saveLocationNameQuery = `
	INSERT INTO locations (id, userID, latitude, longitude, location_name)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (id) DO UPDATE
	SET location_name = $5;
`

func (r LocationRepo) SaveLocationName(ctx context.Context, userID int, locationName string) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id":       userID,
		"location_name": locationName,
	})
	log.Debug("saving the name of the location")

	coord, err := r.GetUserRecentLocation(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "cannot save location name")
	}

	_, err = r.db.ExecContext(ctx, saveLocationNameQuery, coord.LocationID, userID, coord.Latitude, coord.Longitude, locationName)
	if err != nil {
		return errors.Wrap(err, "cannot save location name")
	}

	log.Trace("successfully saved location name")
	return nil
}

const getCoordinatesByCityNameQuery = `
	SELECT latitude, longitude 
	FROM places
	WHERE city = $1;
	`

func (r LocationRepo) GetCoordinatesByCityName(ctx context.Context, locationName string) (bot.Location, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"location_name": locationName,
	})
	log.Debug("getting the coordinates of the location")

	loc := bot.Location{}
	err := r.db.GetContext(ctx, &loc, getCoordinatesByCityNameQuery)
	if err != nil {
		return loc, errors.Wrap(err, "cannot get coordinates by location name")
	}

	return loc, nil
}
