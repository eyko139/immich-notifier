package main

import (
	"github.com/eyko139/immich-notifier/internal/env"
	"net/http"
	"time"
)

func main() {

	env := env.New()

	app := NewApp(env)

	go app.Notifier.StartLoop()

	srv := http.Server{
		ErrorLog:     app.ErrorLog,
		Addr:         ":" + env.AppPort,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
