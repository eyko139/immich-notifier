package main

import (
	"github.com/eyko139/immich-notifier/internal/models"
	"net/http"
)

func (a *App) isAuthenticated(r *http.Request) bool {
	return a.SessionManager.GetBool(r.Context(), "authenticated")
}

func (a *App) GetCurrentUser(r *http.Request) models.User {
	var user models.User
	email := a.SessionManager.GetString(r.Context(), "user_email")
	name := a.SessionManager.GetString(r.Context(), "user_name")
	user.Email = email
	user.Name = name
	return user
}
