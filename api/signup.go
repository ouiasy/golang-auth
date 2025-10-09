package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	apierrors2 "github.com/ouiasy/golang-auth/apierrors"
	"github.com/ouiasy/golang-auth/httputils"

	"github.com/jmoiron/sqlx"
	"github.com/ouiasy/golang-auth/conf"
	"github.com/ouiasy/golang-auth/mailer"
	"github.com/ouiasy/golang-auth/models"
)

type SignupRequest struct {
	UserName string         `json:"username"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Data     map[string]any `json:"data"` // TODO:

	Aud string `json:"-"` // todo
}

type SignupResponse struct {
	ID        uuid.UUID `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *API) Signup(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	params := &SignupRequest{}
	err := httputils.DecodeJSON(r, params)
	if err != nil {
		return apierrors2.BadRequestError(apierrors2.ErrInvalidParameter)
	}

	if err := validateParams(a.Config, params); err != nil {
		return apierrors2.UnprocessableEntityError(apierrors2.ErrValidation)
	}

	user, err := a.Repo.FindUserByEmail(params.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return apierrors2.InternalServerError(apierrors2.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	txx, err := a.Repo.Beginx()
	if err != nil {
		return apierrors2.InternalServerError(apierrors2.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}
	defer txx.Rollback()

	if user != nil {
		if user.IsEmailConfirmed() {
			return apierrors2.BadRequestError(apierrors2.ErrUserAlreadyRegistered)
		}
	} else {
		user, err = registerNewUser(ctx, txx, params)
		if err != nil {
			return apierrors2.InternalServerError(apierrors2.ErrInternalServerError).WithInternalError(
				slog.LevelError, err,
			)
		}
	}

	if !user.IsEmailConfirmed() {
		err := a.EmailClient.SendConfirmationEmail(txx, user, params.Email, a.Config.Mail.SendConfirmationFrequency)
		if err != nil {
			if errors.Is(err, mailer.ErrorMaxFrequencyLimit) {
				return apierrors2.TooManyRequestsError(apierrors2.ErrTooManyRequest)
			}
			return apierrors2.InternalServerError(apierrors2.ErrInternalServerError).WithInternalError(
				slog.LevelError, err,
			)
		}
	}

	if err := txx.Commit(); err != nil {
		return apierrors2.InternalServerError(apierrors2.ErrInternalServerError).WithInternalError(slog.LevelError, err)
	}

	resp := &SignupResponse{
		ID:        user.ID,
		UserName:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return httputils.SendJSON(w, http.StatusOK, resp)
}

// create a new user struct and insert it into db.
func registerNewUser(ctx context.Context, txx *sqlx.Tx, params *SignupRequest) (*models.User, error) {
	user, err := models.NewUser(params.UserName, params.Email, params.Password)
	if err != nil {
		return nil, fmt.Errorf("error while creating user: %s", err)
	}

	query := `INSERT INTO app.users(id, email, username, hashed_password) VALUES ($1, $2, $3, $4)`
	_, err = txx.Exec(query, user.ID, user.Email, user.Username, user.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("error while registering user: %s", err)
	}

	return user, nil
}

var (
	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

func validateParams(config *conf.GlobalConfiguration, params *SignupRequest) error {
	if params.Password == "" || params.Email == "" {
		return fmt.Errorf("require a valid password and email to signup")
	}
	if len(params.Password) < config.App.PasswordMinLength {
		return fmt.Errorf("password should be at least %d characters", config.App.PasswordMinLength)
	}
	if len(params.Password) > config.App.PasswordMaxLength {
		return fmt.Errorf("password should be at most %d characters", config.App.PasswordMaxLength)
	}
	if !emailRegexp.MatchString(params.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}
