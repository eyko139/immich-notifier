package main

import (
	"context"
	"github.com/eyko139/immich-notifier/cmd/util"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/models"
	"github.com/eyko139/immich-notifier/internal/notifier"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
	"time"
)

const (
	NotificationInterval = 5 * time.Second
)

type App struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Helper   *util.Helper
	Immich   *models.ImmichModel
	Users    *models.UserModel
	Notifier *notifier.Notifier
}

func NewApp(env *env.Env) *App {
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	client, err := mongo.Connect(options.Client().ApplyURI(env.DbConnectionString))

	if err != nil {
		errLog.Printf("Failed to connect to DB %s", err)
	}
	pingErr := client.Ping(context.Background(), nil)

	if pingErr != nil {
		errLog.Printf("Database Ping failed, err: %s", pingErr)
	}

	tc, err := models.NewTemplateCache()

	if err != nil {
		errLog.Panicf("Failed to create templateCache, err: %s", err)
	}

	helper := util.New(tc)
	return &App{
		ErrorLog: errLog,
		InfoLog:  infoLog,
		Helper:   helper,
		Users:    models.NewUserModel(client),
		Immich:   models.NewImmichModel(client, env),
		Notifier: notifier.New(client, env, NotificationInterval, models.NewImmichModel(client, env)),
	}
}