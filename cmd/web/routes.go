package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	// create router
	mux := chi.NewRouter()

	// setup middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// define application routes
	mux.Get("/", app.HomePage)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.Logout)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate", app.ActiveAccount)

	mux.Mount("/members", app.authRouter())

	return mux
}

func (app *Config) authRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.Auth)

	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubcribeToPlan)

	return mux
}