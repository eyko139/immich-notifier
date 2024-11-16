package models

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type User struct {
	Name          string              `json:"name" bson:"name"`
	Subscriptions []AlbumSubscription `json:"subscriptions" bson:"subscriptions"`
	Email         string              `json:"email" bson:"email"`
	ID            bson.ObjectID       `json:"id" bson:"_id"`
	ApiKey        string              `json:"apiKey" bson:"apiKey"`
	ChatId        int                 `json:"chat_id" bson:"chatId"`
}

type UserContext struct {
	Email string
	Name  string
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

func (um *UserModel) SaveSubscription(user User) (string, error) {
	filter := bson.D{
		{
			Key: "email", Value: user.Email,
		},
	}
	update := bson.M{"$set": bson.M{"subscriptions": user.Subscriptions}}

	_, err := um.DbClient.Database("Notify").Collection("users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", err
	}
	return "ok", nil
}

func (um *UserModel) FindOrInsertUser(name, email string) (User, error) {
	filter := bson.D{
		{
			Key: "email", Value: email,
		},
	}
	res := um.DbClient.Database("Notify").Collection("users").FindOne(context.TODO(), filter)

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		user := bson.M{
			"email": email,
			"name":  name,
		}
		fmt.Printf("No user found for email: %s, creating...", email)
		_, err := um.DbClient.Database("Notify").Collection("users").InsertOne(context.TODO(), user)
		if err != nil {
			fmt.Println("Error creating user")
		}
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		fmt.Println("failed to decode user")
	}
	return user, nil
}

func (um *UserModel) ActivateSubscriptions(userName string, chatId int) error {

	filter := bson.M{
		"name": userName,
	}

	update := bson.M{"$set": bson.M{"subscriptions.$[].isSubscribed": true, "chatId": chatId}}

	res := um.DbClient.Database("Notify").Collection("users").FindOneAndUpdate(context.TODO(), filter, update)
    if errors.Is(res.Err(), mongo.ErrNoDocuments) {
        fmt.Println("No user found")
    }
	if res.Err() != nil {
        fmt.Println(res.Err())
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		fmt.Println("failed to decode user")
		return nil
	}
	return nil
}
