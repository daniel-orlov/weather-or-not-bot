package utils

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type postgresCfg struct {
	DBName string `mapstructure:"db_name"`
	DBHost string `mapstructure:"db_host"`
	DBPort int    `mapstructure:"db_port"`
	DBUser string `mapstructure:"db_user"`
	DBPass string `mapstructure:"db_pass"`
	DBSSL  bool   `mapstructure:"db_ssl"`
}

func init() {
	pflag.String("db_name", "bot_database", "Database name")
	pflag.String("db_host", "localhost", "Database host")
	pflag.Int("db_port", 5432, "Database port")
	pflag.String("db_user", "root", "Database user")
	pflag.String("db_pass", "1234", "Database password")
	pflag.Bool("db_ssl", false, "Is database SSL mode on")

	sql.Register("instrumented-postgres", stdlib.GetDefaultDriver())
}

// NewDBFromEnv establishes a new db connection and returns a wrapper.
func NewDBFromEnv() *sqlx.DB {
	cfg := getDataFromEnv()
	logrus.WithFields(logrus.Fields{
		"db_user": cfg.DBUser,
		"db_host": cfg.DBHost,
		"db_port": cfg.DBPort,
		"db_name": cfg.DBName,
	}).Info("Establishing a new database connection")

	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName))

	if !cfg.DBSSL {
		buf.WriteString(" sslmode=disable")
	}

	connStr := buf.String()

	db, err := sql.Open("instrumented-postgres", connStr)
	if err != nil {
		logrus.WithError(err).Panic("Cannot open driver with connection string")
	}

	dbx := sqlx.NewDb(db, "postgres")

	if err := dbx.Ping(); err != nil {
		logrus.WithError(err).Fatal("Cannot ping database")
	}

	return dbx
}

func getDataFromEnv() *postgresCfg {
	var cfg postgresCfg

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Cannot get db conn cfg from envs")
	}

	return &cfg
}
