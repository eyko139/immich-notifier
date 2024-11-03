package main

import "net/http"

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
		a.InfoLog.Println(apiKey)

		albums, _ := a.Immich.FetchAlbums(apiKey)
		templateData := a.Helper.NewTemplateData(albums, apiKey)

		a.Helper.ReturnHtml(w, "albums.html", templateData)
	}
}

func (a *App) notifyPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type FormInput struct {
			album []string
			email string
			apiKey string
		}
		formValues := &FormInput{
			album: []string{},
		}
		err := r.ParseForm()
		if err != nil {
			a.ErrorLog.Println("Failed to parse form")
		}
		for key, value := range r.Form {
			if key == "album" {
				for _, val := range value {
					formValues.album = append(formValues.album, val)
				}
			}
			formValues.email = r.Form["email"][0]
			formValues.apiKey = r.Form["apiKey"][0]
		}
		res, _ := a.Users.SaveSubscription("test", formValues.album, formValues.apiKey)
		a.InfoLog.Println(res)
	}
}
