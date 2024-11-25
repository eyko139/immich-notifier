package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/eyko139/immich-notifier/cmd/util"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/models"
	"github.com/eyko139/immich-notifier/internal/models/mocks"
	"github.com/eyko139/scs/v2"
)


func newTestApplication(env *env.Env) *App {

    var sessionManager *scs.SessionManager

	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	tc, err := models.NewTemplateCache()

	if err != nil {
		errLog.Panicf("Failed to create templateCache, err: %s", err)
	}

	helper := util.New(tc, errLog, infoLog, "test")

    return &App{
        ErrorLog: errLog, 
        InfoLog: infoLog,
        Helper: helper,
        Users: &mocks.UserModel{},
        Immich: &mocks.ImmichModel{},
        SessionManager: sessionManager,
        Env: env,
    }
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, handler http.Handler) *testServer {
	ts := httptest.NewServer(handler)
	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar
	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx
	// response is received by the client, and by always returning a
	// http.ErrUseLastResponse error it forces the client to immediately return
	// the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

    return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, url string) (int, http.Header, string) {

    rs, err := ts.Client().Get(ts.URL + url)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func MockAuthenticationMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}

