package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (a *App) Routes() http.Handler {

	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("../../ui/static/"))

	dynamic := alice.New(a.SessionManager.LoadAndSave)
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.Handler(http.MethodPost, "/bothook", dynamic.ThenFunc(a.botHook()))
	router.Handler(http.MethodGet, "/login", dynamic.ThenFunc(a.login()))
	router.Handler(http.MethodGet, "/callback", dynamic.ThenFunc(a.handleCallback()))

	protected := dynamic.Append(a.requireAuthentication)
	router.Handler(http.MethodGet, "/", protected.ThenFunc(a.home()))
	router.Handler(http.MethodPost, "/apikey", protected.ThenFunc(a.apiKeyPost()))
	router.Handler(http.MethodPost, "/notifypost", protected.ThenFunc(a.notifyPost()))
	standard := alice.New(a.LogRequests, secureHeaders)
	return standard.Then(router)

}
