package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bot "gopkg.in/telegram-bot-api.v4"
	"weather-or-not-bot/internal/types"
)

type UserLocationRepo struct {
	db *sqlx.DB
}

func NewUserLocationRepo(db *sqlx.DB) *UserLocationRepo {
	return &UserLocationRepo{db: db}
}

const addUserLocationByCoordinatesQuery = `
	INSERT INTO locations (user_id, latitude, longitude)
		VALUES ($1, $2, $3);
`

func (r *UserLocationRepo) AddUserLocationByCoordinates(ctx context.Context, userID int, loc *bot.Location) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id": userID,
		"lat":     loc.Latitude,
		"long":    loc.Longitude,
	})
	log.Debug("Adding the location by coordinates")

	_, err := r.db.ExecContext(ctx, addUserLocationByCoordinatesQuery, userID,
		fmt.Sprintf("%f", loc.Latitude),
		fmt.Sprintf("%f", loc.Longitude),
	)
	if err != nil {
		return errors.Wrap(err, "cannot add location by coordinates")
	}

	log.Debug("Successfully added location by coordinates")

	return nil
}

const getUserRecentLocationQuery = `
	SELECT id, latitude, longitude
	FROM locations
	WHERE user_id = $1
	ORDER BY id DESC
	LIMIT 1;
`

func (r *UserLocationRepo) GetUserRecentLocation(ctx context.Context, userID int) (*types.UserCoordinates, error) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id": userID,
	})
	log.Debug("getting user's recent coordinates from db")

	userLocation := types.UserCoordinates{}
	err := r.db.GetContext(ctx, &userLocation, getUserRecentLocationQuery, userID)
	if err != nil {
		return &userLocation, errors.Wrap(err, "cannot get user's recent location")
	}

	return &userLocation, nil
}

const saveLocationNameQuery = `
	UPDATE locations
	SET location_name = $1
	WHERE id = $2;
`

func (r *UserLocationRepo) SaveUserLocationName(ctx context.Context, userID int, locationName string) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id":       userID,
		"location_name": locationName,
	})
	log.Debug("saving the name of the location")

	// TODO get rid of two queries
	coord, err := r.GetUserRecentLocation(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "cannot save location name")
	}

	_, err = r.db.ExecContext(ctx, saveLocationNameQuery, locationName, coord.LocationID)
	if err != nil {
		return errors.Wrap(err, "cannot save location name")
	}

	log.Debug("successfully saved the location name")
	return nil
}
