package api

import (
	"net/http"

	"github.com/ouiasy/golang-auth/conf"
	"github.com/ouiasy/golang-auth/mailer"
	"github.com/ouiasy/golang-auth/repository"
)

type API struct {
	Config  *conf.GlobalConfiguration
	Handler http.Handler
	Repo    *repository.Repository

	EmailClient *mailer.EmailClient
}

func NewApi(globalConfig *conf.GlobalConfiguration, repo *repository.Repository, eClient *mailer.EmailClient) *API {
	api := &API{
		Config:      globalConfig,
		Repo:        repo,
		EmailClient: eClient,
	}

	r := newRouter()

	r.Route("/", func(r *router) {
		r.Post("/signup", api.Signup)
	})

	// todo: add handlers

	api.Handler = r

	return api
}
