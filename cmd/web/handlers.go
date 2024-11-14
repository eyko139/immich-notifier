package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"github.com/eyko139/immich-notifier/internal/models"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"time"
)

var stateStore = make(map[string]bool)

func (a *App) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mail := a.SessionManager.GetString(r.Context(), "user_email")
		userContext := models.UserContext{Email: mail}
		a.Helper.Render(w, "home.html", &userContext)
	}
}

func (a *App) apiKeyPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			a.ErrorLog.Printf("Error parsing form, err: %s", err)
		}
		apiKey := r.Form.Get("apiKey")

		albums, _ := a.Immich.FetchAlbums(apiKey)
		user, err := a.Users.FindUser(apiKey)
		if err != nil {
			a.ErrorLog.Println("no user found")
		}
		for _, sub := range user.Subscriptions {
			for idx, album := range albums {
				if sub.Id == album.Id {
					albums[idx].IsSubscribed = true
				}
			}
		}
		for _, album := range albums {
			a.Immich.InsertOrAlbum(album)
		}
		templateData := a.Helper.NewTemplateData(albums, apiKey)
		a.Helper.ReturnHtml(w, "albums.html", templateData)
	}
}

func (a *App) notifyPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		var user models.User
		user.ApiKey = r.Form["apiKey"][0]
		user.Subscriptions = []models.AlbumSubscription{}

		if err != nil {
			a.ErrorLog.Println("Failed to parse form")
		}
		for key, value := range r.Form {
			if key == "album" {
				for _, val := range value {
					var subscription models.AlbumSubscription
					album, err := a.Immich.FetchAlbumsDetails(val, user.ApiKey)
					if err != nil {
						a.ErrorLog.Printf("Error fetching api details: %s", err)
					}
					subscription.Id = album.Id
					subscription.AlbumName = album.AlbumName
					subscription.LastNotified = time.Now()
					subscription.IsSubscribed = false
					user.Subscriptions = append(user.Subscriptions, subscription)
				}
			}
		}
		res, _ := a.Users.SaveSubscription(user)
		a.InfoLog.Println(res)

		a.Helper.ReturnPlainHtml(w, "notify.html", struct{ ApiKey string }{ApiKey: user.ApiKey})
	}
}

func (a *App) botHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var botResponse models.BotResponse
		bytes, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(bytes, &botResponse)
		if err != nil {
			a.ErrorLog.Printf("Error parsing bot response: %s", err)
		}
		apiKey := botResponse.Message.Text[7:]
		err = a.Users.ActivateSubscriptions(apiKey, botResponse.Message.From.Id)
		if err != nil {
			a.ErrorLog.Println(err)
		}
	}
}

func (a *App) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, _ := generateState()
		stateStore[state] = true
		url := a.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func (a *App) handleCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		// Verify state parameter
		state := r.URL.Query().Get("state")
		if !stateStore[state] {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}
		delete(stateStore, state)

		// Exchange the authorization code for a token
		token, err := a.OauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Extract ID Token from OAuth2 token
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token in token response", http.StatusInternalServerError)
			return
		}

		// Parse and verify ID Token
		verifier := a.OauthProvider.Verifier(&oidc.Config{ClientID: a.Env.OidcClientId})
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get user information from ID token
		var claims struct {
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse ID Token claims: "+err.Error(), http.StatusInternalServerError)
			return
		}
		sessionManager.Put(r.Context(), "authenticated", true)
		sessionManager.Put(r.Context(), "user_email", claims.Email)
		sessionManager.Put(r.Context(), "user_name", claims.Name)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
func generateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
