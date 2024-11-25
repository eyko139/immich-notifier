package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eyko139/immich-notifier/internal/assert"
	"github.com/eyko139/immich-notifier/internal/env"
	_ "github.com/eyko139/immich-notifier/internal/test_utils"
)

func TestHealth(t *testing.T) {
	env := env.New()
	app := newTestApplication(env)

	ts := newTestServer(t, app.Routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/health")

	assert.AssertEqual(t, code, http.StatusOK)
	assert.AssertEqual(t, body, "OK")
}

func TestHome(t *testing.T) {

    t.Setenv("APP_ENV", "test")

	env := env.New()
	appWithUser := newTestApplication(env)

	ts := newTestServer(t, appWithUser.Routes())

	defer ts.Close()

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
		app        *App
		userName   string
	}{
		{
			name:       "Fetch single album",
			url:        "/",
			wantStatus: http.StatusOK,
			wantBody:   "mockAlbum",
			app:        appWithUser,
			userName:   "nobody",
		},
		{
			name:       "Toggle on; album is subbed",
			url:        "/",
			wantStatus: http.StatusOK,
			wantBody:   "checked",
			app:        appWithUser,
			userName:   "active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

            req := httptest.NewRequest("GET", ts.URL + tt.url, nil)
            rr := httptest.NewRecorder()

            ctx, _ := tt.app.SessionManager.Load(req.Context(), tt.name)

            tt.app.SessionManager.Put(ctx, "authenticated", true)
            tt.app.SessionManager.Put(ctx, "user_name", tt.userName)
            tt.app.SessionManager.Put(ctx, "user_email", "test@testmail.com")

            appWithUser.Routes().ServeHTTP(rr, req.WithContext(ctx))
            by, _ := io.ReadAll(rr.Body)

			assert.AssertEqual(t, rr.Code, tt.wantStatus)
			assert.StringContains(t, string(by), tt.wantBody)
		})
	}
}
