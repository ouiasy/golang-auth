package repository

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ouiasy/golang-auth/conf"
)

type Repository struct {
	*sqlx.DB
}

func dialPostgres(config *conf.GlobalConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", config.DB.Dsn())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(config *conf.GlobalConfiguration) (*Repository, error) {
	db, err := dialPostgres(config)
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}
