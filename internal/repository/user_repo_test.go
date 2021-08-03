package repository

import (
	"context"
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	bot "gopkg.in/telegram-bot-api.v4"
)

func TestUserDataRepo_AddUserIfNotExists(t *testing.T) {
	ctx := context.Background()
	user := &bot.User{
		ID:           123,
		FirstName:    "John",
		LastName:     "Doe",
		UserName:     "the_john",
		LanguageCode: "en",
		IsBot:        false,
	}

	const expectedQuery = `
	INSERT INTO users (user_id, username, first_name, last_name, language_code, is_bot)
	VALUES (.+)
	ON CONFLICT (user_id) DO NOTHING;
`

	tests := []struct {
		name    string
		prepare func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			"1. Error on add user",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(user.ID, user.UserName, user.FirstName, user.LastName, user.LanguageCode, user.IsBot).
					WillReturnError(errors.New("some error"))
			},
			true,
		},
		{
			"2. SQL No rows - user exists",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(user.ID, user.UserName, user.FirstName, user.LastName, user.LanguageCode, user.IsBot).
					WillReturnError(sql.ErrNoRows)
			},
			false,
		},
		{
			"3. Success on add user",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(expectedQuery).
					WithArgs(user.ID, user.UserName, user.FirstName, user.LastName, user.LanguageCode, user.IsBot).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			defer func() {
				if expErr := mock.ExpectationsWereMet(); expErr != nil {
					t.Errorf("UserDataRepo.AddUserIfNotExists() there were unfulfilled expectations: %s", expErr)
				}
			}()

			tt.prepare(mock)

			repo := NewUserDataRepo(sqlx.NewDb(db, "postgres"))
			err = repo.AddUserIfNotExists(ctx, user)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDataRepo.AddUserIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
