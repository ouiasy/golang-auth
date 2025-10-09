package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ouiasy/golang-auth/models"
)

func (r *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := `SELECT * FROM app.users WHERE email = $1`

	err := r.Get(&user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("error while querying user: %s", err)
	}

	return &user, nil
}

func FindUserByToken(db sqlx.Ext, token string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM app.users WHERE confirmation_token = $1`
	err := sqlx.Get(db, user, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func ConfirmUser(db sqlx.Ext, user *models.User) error {
	now := time.Now()
	query := `UPDATE app.users SET email_confirmed_at = $1, confirmation_token = NULL WHERE id = $2`
	_, err := db.Exec(query, now, user.ID)
	if err != nil {
		return err
	}

	return nil
}
