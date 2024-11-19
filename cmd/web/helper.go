package main

import (
	"github.com/eyko139/immich-notifier/internal/models"
	customErrors "github.com/eyko139/immich-notifier/internal/errors"
	"net/http"
)

func (a *App) isAuthenticated(r *http.Request) bool {
	return a.SessionManager.GetBool(r.Context(), "authenticated")
}

func (a *App) GetCurrentSessionUser(r *http.Request) (*models.User, error) {
	var user models.User
	email := a.SessionManager.GetString(r.Context(), "user_email")
	name := a.SessionManager.GetString(r.Context(), "user_name")
	chatId := a.SessionManager.GetInt(r.Context(), "user_chatId")

    if email == "" || name == "" {
        return nil, &customErrors.NoUserInSessionError{Session: a.SessionManager.Token(r.Context()), Message: "session"}
    }

	user.Email = email
	user.Name = name
    user.ChatId = chatId
	return &user, nil
}
