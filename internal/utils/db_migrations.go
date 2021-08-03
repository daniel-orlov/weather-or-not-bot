package utils

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// RunInitMigration prepares a db instance by running initial migrations.
func RunInitMigration(ctx context.Context, db *sqlx.DB) error {
	log := ctxlogrus.Extract(ctx)
	log.Info("Running init migrations")

	_, err := db.ExecContext(ctx, initMigration)
	if err != nil {
		return errors.Wrap(err, "cannot run init migrations")
	}

	return nil
}

const initMigration = `
	CREATE TABLE IF NOT EXISTS users
(
	user_id BIGINT PRIMARY KEY,
	username VARCHAR(64) NOT NULL DEFAULT '',
	first_name VARCHAR(64) NOT NULL DEFAULT '',
	last_name VARCHAR(64) NOT NULL DEFAULT '',
	language_code VARCHAR(8) NOT NULL DEFAULT 'en',
	is_bot BOOLEAN NOT NULL DEFAULT false,
	CONSTRAINT unique_user_id UNIQUE (user_id)
);

	CREATE TABLE IF NOT EXISTS locations
(
	id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id BIGINT NOT NULL,
	latitude VARCHAR(64) NOT NULL DEFAULT '',
	longitude VARCHAR(64) NOT NULL DEFAULT '',
	location_name VARCHAR(64) NOT NULL DEFAULT ''
);
	drop table if exists world_cities;

	CREATE TABLE IF NOT EXISTS world_cities
(
	city VARCHAR(64) NOT NULL DEFAULT '',
	city_ascii VARCHAR(64) NOT NULL DEFAULT '',
	lat VARCHAR(64) NOT NULL DEFAULT '',
	long VARCHAR(64) NOT NULL DEFAULT '',
	country VARCHAR(64) NOT NULL DEFAULT '',
	iso2 VARCHAR(2) NOT NULL DEFAULT '',
	iso3 VARCHAR(3) NOT NULL DEFAULT ''
);

	COPY world_cities
    FROM '/world_cities.csv'
    DELIMITER ','
    CSV HEADER;
`
