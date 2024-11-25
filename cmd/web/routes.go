package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
	"os"
	"path/filepath"
)

func (a *App) Routes() http.Handler {

	router := httprouter.New()

	cwd, err := os.Getwd() // Get the current working directory
	if err != nil {
		panic(err)
	}
	staticPath := filepath.Join(cwd, "ui/static")

	fileServer := http.FileServer(http.Dir(staticPath))

	dynamic := alice.New(a.SessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.Handler(http.MethodPost, "/bothook", dynamic.ThenFunc(a.botHook()))
	router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(a.login()))
	router.Handler(http.MethodGet, "/callback", dynamic.ThenFunc(a.handleCallback()))
    router.Handler(http.MethodGet, "/logout-success", dynamic.ThenFunc(a.logoutSuccess()))
    router.Handler(http.MethodGet, "/health", dynamic.ThenFunc(a.health()))

	protected := dynamic.Append(a.requireAuthentication)
	router.Handler(http.MethodGet, "/", protected.ThenFunc(a.home()))
    router.Handler(http.MethodPost, "/subscribe/:albumId", protected.ThenFunc(a.subAlbumPost()))
    router.Handler(http.MethodGet, "/logout", protected.ThenFunc(a.logout()))
	standard := alice.New(a.LogRequests, secureHeaders)
	return standard.Then(router)
}
