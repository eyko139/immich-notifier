package main

import "net/http"

func (a *App) isAuthenticated(r *http.Request) bool {
	return a.SessionManager.GetBool(r.Context(), "authenticated")
}
