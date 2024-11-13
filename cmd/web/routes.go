package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (a *App) Routes() http.Handler {

	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("../../ui/static/"))

	dynamic := alice.New()
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(a.home()))
	router.Handler(http.MethodPost, "/bothook", dynamic.ThenFunc(a.botHook()))
	router.Handler(http.MethodPost, "/apikey", dynamic.ThenFunc(a.apiKeyPost()))
	router.Handler(http.MethodPost, "/notifypost", dynamic.ThenFunc(a.notifyPost()))
	standard := alice.New(a.LogRequests, secureHeaders)
	return standard.Then(router)

}
