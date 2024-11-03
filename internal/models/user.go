package models

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type User struct {
	Name          string              `json:"Name" bson:"name"`
	Subscriptions []AlbumSubscription `json:"Subscriptions" bson:"Subscriptions"`
	Email         string              `json:"email" bson:"email"`
	ID            bson.ObjectID       `json:"id" bson:"_id"`
	ApiKey        string              `json:"apiKey" bson:"apiKey"`
}

type AlbumSubscription struct {
	AlbumName    string    `json:"albumName" bson:"albumName"`
	Id           string    `json:"id" bson:"id"`
	LastNotified time.Time `json:"lastNotified" bson:"lastNotified"`
	IsSubscribed bool      `json:"isSubscribed" bson:"isSubscribed"`
}

type UserModel struct {
	DbClient *mongo.Client
}

func NewUserModel(client *mongo.Client) *UserModel {
	return &UserModel{
		DbClient: client,
	}
}

func (um *UserModel) Insert(name string, subscribedAlbums []string, apiKey string) (string, error) {
	return "", nil
}

func (um *UserModel) SaveSubscription(user User) (string, error) {
	_, err := um.DbClient.Database("Notify").Collection("users").InsertOne(context.TODO(), user, nil)
	if err != nil {
		return "", err
	}
	return "ok", nil
}
