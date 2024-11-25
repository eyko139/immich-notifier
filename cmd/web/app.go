package main

import (
	"context"
	"github.com/eyko139/scs/mongodbstore"
	"github.com/eyko139/scs/v2"
	"github.com/coreos/go-oidc"
	"github.com/eyko139/immich-notifier/cmd/util"
	"github.com/eyko139/immich-notifier/internal/auth"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/models"
	"github.com/eyko139/immich-notifier/internal/notifier"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
)

type App struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Helper         *util.Helper
	Immich         models.ImmichModelInterface
	Users          models.UserModelInterface
	Notifier       *notifier.Notifier
	OauthConfig    *oauth2.Config
	OauthProvider  *oidc.Provider
	SessionManager *scs.SessionManager
	Env            *env.Env
}

var sessionManager *scs.SessionManager

func NewApp(env *env.Env) *App {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)


    infoLog.Printf("Connecting to DB %s", env.DbConnectionString)
    options := options.Client().SetTimeout(5 * time.Second).ApplyURI(env.DbConnectionString)

	client, err := mongo.Connect(options)
	if err != nil {
		errLog.Printf("Failed to connect to DB %s", err)
	}

    pingCtx := context.Background()
    ctx, cancel := context.WithTimeout(pingCtx, 5 * time.Second)
    defer cancel()
	pingErr := client.Ping(ctx, nil)

	if pingErr != nil {
		errLog.Printf("Database Ping failed, err: %s", pingErr)
	}

    db:= client.Database("Notify")
    
	sessionManager = scs.New()
    sessionManager.Store = mongodbstore.New(db)
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Secure = false // Enable secure cookie in production




	tc, err := models.NewTemplateCache()

	if err != nil {
		errLog.Panicf("Failed to create templateCache, err: %s", err)
	}

	oAuthConfig, provider := auth.NewOauthConfig(env.OidcIssuerUrl, env.OidcClientId, env.OidcClientSecret, env.OidcRedirectUrl)

	helper := util.New(tc, errLog, infoLog, env.AppVersion)

	return &App{
		ErrorLog:       errLog,
		InfoLog:        infoLog,
		Helper:         helper,
		Users:          models.NewUserModel(db),
		Immich:         models.NewImmichModel(db, env),
		Notifier:       notifier.New(client, env, models.NewImmichModel(db, env), errLog, infoLog),
		OauthConfig:    oAuthConfig,
		OauthProvider:  provider,
		Env:            env,
		SessionManager: sessionManager,
	}
}
