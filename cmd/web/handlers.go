package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/eyko139/immich-notifier/internal/models"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

var stateStore = make(map[string]bool)

func (a *App) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mail := a.SessionManager.GetString(r.Context(), "user_email")
		name := a.SessionManager.GetString(r.Context(), "user_name")

		albums, _ := a.Immich.FetchAlbums(mail)
		user, err := a.Users.FindOrInsertUser(name, mail)

		a.SessionManager.Put(r.Context(), "user_chatId", user.ChatId)

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
		templateData := a.Helper.NewTemplateData(albums, mail, name, user.ChatId != 0, user.ID.Hex())
		a.Helper.Render(w, "home.html", templateData)
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
				userId := parts[1]
				if err := a.Users.ActivateSubscriptions(userId, botResponse.Message.From.Id); err != nil {
					a.ErrorLog.Println("Failed to activate subscription, error: " + err.Error())
					return
				}
				a.Notifier.SendTelegramMessage(botResponse.Message.From.Id, fmt.Sprintf("Bot activated, return to website: %s", a.Env.WebsiteURL))
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
		if claims.Email == "" || claims.Name == "" {
			http.Error(w, "Missing email or username in token: "+err.Error(), http.StatusInternalServerError)
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

func (a *App) subAlbumPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id := params.ByName("albumId")
		a.InfoLog.Println("Subscribing to album: " + id)

		user := a.GetCurrentSessionUser(r)

		var subscription models.AlbumSubscription
		album, err := a.Immich.FetchAlbumsDetails(id)
		if err != nil {
			a.ErrorLog.Printf("Error fetching api details: %s", err)
		}
		subscription.Id = album.Id
		subscription.AlbumName = album.AlbumName
		subscription.LastNotified = time.Now()
		user.Subscriptions = append(user.Subscriptions, subscription)

		if err := a.Users.UpdateSubscription(user.Email, subscription); err != nil {
			a.ErrorLog.Printf("Failed to update album subscription: %s", err.Error())
		}
	}
}

func (a *App) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionManager.Destroy(r.Context())
		w.Header().Set("HX-Location", "/logout-success")
	}
}

func (a *App) logoutSuccess() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.Helper.Render(w, "logout.html", nil)
	}
}
