package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"github.com/eyko139/immich-notifier/internal/models"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var stateStore = make(map[string]bool)

func (a *App) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mail := a.SessionManager.GetString(r.Context(), "user_email")
		name := a.SessionManager.GetString(r.Context(), "user_name")
		albums, _ := a.Immich.FetchAlbums()
		user, err := a.Users.FindOrInsertUser(name, mail)
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
		templateData := a.Helper.NewTemplateData(albums, mail, name)
		a.Helper.Render(w, "home.html", templateData)
	}
}

func (a *App) notifyPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		user := a.GetCurrentUser(r)

		user.Subscriptions = []models.AlbumSubscription{}

		if err != nil {
			a.ErrorLog.Println("Failed to parse form")
		}
		for key, value := range r.Form {
			if key == "album" {
				for _, val := range value {
					var subscription models.AlbumSubscription
					album, err := a.Immich.FetchAlbumsDetails(val)
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

		type NotifyEmail struct {
			UserName string
		}

		data := NotifyEmail{UserName: url.QueryEscape(user.Name)}

		a.Helper.ReturnPlainHtml(w, "notify.html", data)
	}
}

func (a *App) botHook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var botResponse models.BotResponse

		if err := json.NewDecoder(r.Body).Decode(&botResponse); err != nil {
			a.ErrorLog.Printf("Error parsing bot response: %s", err)
		}

		if strings.HasPrefix(botResponse.Message.Text, "/start") {
			parts := strings.SplitN(botResponse.Message.Text, " ", 2)
			if len(parts) > 1 {
				a.InfoLog.Println("Bothook query: " + parts[1])
				userName := parts[1]
				err := a.Users.ActivateSubscriptions(userName, botResponse.Message.From.Id)
				if err != nil {
					a.ErrorLog.Println("Failed to activate subscription")
				}
			} else {
				a.InfoLog.Println("Bothook called with no parameters")
			}
		}
	}
}

func (a *App) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authCodeURL string
		state, _ := generateState()
		stateStore[state] = true
		if a.Env.AppEnv == "development" {
			authCodeURL = "/"
		} else {
			authCodeURL = a.OauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
		}

		http.Redirect(w, r, authCodeURL, http.StatusFound)
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

		user, _ := a.Users.FindOrInsertUser(claims.Name, claims.Email)

		a.InfoLog.Printf("created user: %+v", user)

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
