package main

import (
	"github.com/eyko139/immich-notifier/internal/models"
	"net/http"
	"time"
)

func (a *App) home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.Helper.Render(w, "home.html", nil)
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
		user.Email = r.Form["email"][0]
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
					subscription.IsSubscribed = true
					user.Subscriptions = append(user.Subscriptions, subscription)
				}
			}
		}
		res, _ := a.Users.SaveSubscription(user)
		a.InfoLog.Println(res)
	}
}
