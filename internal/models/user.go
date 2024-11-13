package models

import (
	"context"
	"fmt"
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
	ChatId        int                 `json:"chat_id" bson:"chatId"`
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

func (um *UserModel) FindUser(apiKey string) (User, error) {
	filter := bson.D{
		{
			"apiKey", apiKey,
		},
	}
	res := um.DbClient.Database("Notify").Collection("users").FindOne(context.TODO(), filter)
	if res.Err() != nil {
		fmt.Println("No user found for apikey")
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		fmt.Println("failed to decode user")
	}
	return user, nil
}

func (um *UserModel) ActivateSubscriptions(apiKey string, chatId int) error {
	filter := bson.M{
		"apiKey": apiKey,
	}

	update := bson.M{"$set": bson.M{"Subscriptions.$[].isSubscribed": true, "chatId": chatId}}

	res := um.DbClient.Database("Notify").Collection("users").FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		fmt.Println("No user found for apikey")
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		fmt.Println("failed to decode user")
		return nil
	}
	return nil
}
