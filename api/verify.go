package api

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/ouiasy/golang-auth/apierrors"
	"github.com/ouiasy/golang-auth/repository"
)

func (a *API) Verify(w http.ResponseWriter, r *http.Request) error {
	// todo: add redirect param and sanitize it
	token := r.URL.Query().Get("token")

	tx, err := a.Repo.Beginx()
	defer tx.Rollback()
	if err != nil {
		return apierrors.InternalServerError(apierrors.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	user, err := repository.FindUserByToken(tx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apierrors.BadRequestError(apierrors.ErrUserNotFound)
		}
		return apierrors.InternalServerError(apierrors.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	tokenExpiration := user.ConfirmationSentAt.Add(a.Config.App.ConfirmationTokenExpiration)
	if time.Now().After(tokenExpiration) {
		return apierrors.BadRequestError(apierrors.ErrConfirmationTokenExpired)
	}

	if err := repository.ConfirmUser(tx, user); err != nil {
		return apierrors.InternalServerError(apierrors.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	// todo: issue token

	if err := tx.Commit(); err != nil {
		return apierrors.InternalServerError(apierrors.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	return nil
}
