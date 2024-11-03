package models

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type User struct {
	Name             string             `json:"Name" bson:"name"`
	SubscribedAlbums []string           `json:"SubscribedAlbums" bson:"albums"`
	Email            string             `json:"email" bson:"email"`
	ApiKey           string             `json:"apiKey" bson:"apiKey"`
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

func (um *UserModel) SaveSubscription(email string, subscribedAlbums []string, apiKey string) (string, error) {
	_, err := um.DbClient.Database("Notify").Collection("users").InsertOne(context.TODO(), bson.D{
		{Key: "email", Value: email},
		{Key: "albums", Value: subscribedAlbums},
		{Key: "apiKey", Value: apiKey},
	})
	if err != nil {
		return "", err
	}
	return "ok", nil
}
