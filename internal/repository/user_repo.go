package repository

import (
	"context"
	"database/sql"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bot "gopkg.in/telegram-bot-api.v4"
)

type UserDataRepo struct {
	db *sqlx.DB
}

func NewUserDataRepo(db *sqlx.DB) *UserDataRepo {
	return &UserDataRepo{db: db}
}

const addUserIfNotExistsQuery = `
	INSERT INTO users (user_id, username, first_name, last_name, language_code, is_bot)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (user_id) DO NOTHING;
`

func (r UserDataRepo) AddUserIfNotExists(ctx context.Context, user *bot.User) error {
	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"username": user.UserName,
		"user_id":  user.ID,
	})
	log.Debug("Adding user to db")

	_, err := r.db.ExecContext(ctx, addUserIfNotExistsQuery, user.ID, user.UserName, user.FirstName, user.LastName, user.LanguageCode, user.IsBot)
	if errors.Is(err, sql.ErrNoRows) {
		log.Debug("User already exists")
		return nil
	}

	if err != nil {
		return errors.Wrap(err, "cannot add user")
	}

	return nil
}
