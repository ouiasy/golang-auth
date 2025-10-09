package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               uuid.UUID  `db:"id"`
	Role             string     `db:"role"`
	Email            string     `db:"email"`
	Username         string     `db:"username"`
	HashedPassword   string     `db:"hashed_password"`
	EmailConfirmedAt *time.Time `db:"email_confirmed_at"`
	LastSignInAt     *time.Time `db:"last_sign_in_at"`

	ConfirmationToken  string     `db:"confirmation_token"`
	ConfirmationSentAt *time.Time `db:"confirmation_sent_at"`

	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// returns a new User struct with (ID, Role, Email, Username, HasshedPw)
func NewUser(username, email, password string) (*User, error) {
	uid := uuid.New()

	pw, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:             uid,
		Role:           "user",
		Email:          strings.ToLower(email),
		Username:       username,
		HashedPassword: pw,
	}

	return user, nil
}

// check EmailConfirmedAt field (bool)
func (u *User) IsEmailConfirmed() bool {
	return u.EmailConfirmedAt != nil
}

func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pw), nil
}
