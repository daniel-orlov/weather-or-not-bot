package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"weather-or-not-bot/internal/types"
)

type UserLocationRepo struct {
	db *sqlx.DB
}

func NewUserLocationRepo(db *sqlx.DB) *UserLocationRepo {
	return &UserLocationRepo{db: db}
}

const addUserLocationByCoordinatesQuery = `
	INSERT INTO locations (userID, latitude, longitude)
	VALUES ($1, $2, $3);
`

func (r UserLocationRepo) AddUserLocationByCoordinates(ctx context.Context, userID int, lat, long float64) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id": userID,
		"lat":     lat,
		"long":    long,
	})
	log.Debug("adding location to db")

	_, err := r.db.ExecContext(ctx, addUserLocationByCoordinatesQuery, userID,
		fmt.Sprintf("%v", lat),
		fmt.Sprintf("%v", long),
	)
	if err != nil {
		return errors.Wrap(err, "cannot add location by coordinates")
	}

	return nil
}

const getUserRecentLocationQuery = `
	SELECT id, latitude, longitude
	FROM locations
	WHERE user_id = $1
	ORDER BY id DESC
	LIMIT 1;
`

func (r UserLocationRepo) GetUserRecentLocation(ctx context.Context, userID int) (*types.UserCoordinates, error) {
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
