package repository

import (
	"context"
	"fmt"
	"github.com/daniel-orlov/weather-or-not-bot/internal/types"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bot "gopkg.in/telegram-bot-api.v4"
)

type CoordinatesRepo struct {
	db *sqlx.DB
}

func NewCoordinatesRepo(db *sqlx.DB) *CoordinatesRepo {
	return &CoordinatesRepo{db: db}
}

const addUserLocationByCoordinatesQuery = `
	INSERT INTO locations (userID, latitude, longitude)
	VALUES ($1, $2, $3);
`

func (r CoordinatesRepo) AddUserLocationByCoordinates(ctx context.Context, userID int, lat, long float64) error {
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

func (r CoordinatesRepo) GetUserRecentLocation(ctx context.Context, userID int) {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"user_id": userID,
	})
	log.Debug("getting user's recent coordinates from db")

	userLocation := types.UserCoordinates{}
	err := r.db.GetContext(ctx)
}

func getMostRecentLocation(conn *pgx.Conn, userId int) (bot.Location, int) {

	loc := bot.Location{}

	var lat, long string
	var id int
	err := conn.QueryRow(context.Background(), sqlQuery).Scan(&id, &lat, &long)
	if err != nil {
		err = errors.Wrap(err, "FAILED: QueryRow when adding LocationByCoords")
		fmt.Println(err)
	}
	fmt.Println("LONG", long, "LAT", lat)
	loc.Latitude, err = strconv.ParseFloat(lat, 64)
	if err != nil {
		err = errors.Wrap(err, "Latitude parsing failed")
		fmt.Println(err)
	}
	loc.Longitude, err = strconv.ParseFloat(long, 64)
	if err != nil {
		err = errors.Wrap(err, "Longitude parsing failed")
		fmt.Println(err)
	}
	return loc, id
}

func nameMostRecentLocation(name string, conn *pgx.Conn, userId int) {

	fmt.Println("EXECUTING: nameMostRecentLocation")
	loc, id := getMostRecentLocation(conn, userId)
	sqlQuery := fmt.Sprintf(
		`INSERT INTO locations (id, "user", latitude, longitude, location)
				VALUES (%v, %v, '%v', '%v', '%v')
				ON CONFLICT (id) DO UPDATE
				SET location = '%v'
				`,
		id, userId, loc.Latitude, loc.Longitude, name, name,
	)
	fmt.Println(sqlQuery)
	err := conn.QueryRow(context.Background(), sqlQuery).Scan()
	if err.Error() == "no rows in result set" { //I have doubts about this one. Not sure if this is right
		fmt.Println("SUCCEEDED: QueryRow when naming MostRecentLocation")
	} else if err != nil {
		err = errors.Wrap(err, "FAILED: QueryRow when naming MostRecentLocation")
		fmt.Println(err)
	}
}
