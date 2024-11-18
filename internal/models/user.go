package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
	Email             string
	Name              string
	TelegramAvailable bool
}

type AlbumSubscription struct {
	AlbumName    string    `json:"albumName" bson:"albumName"`
	Id           string    `json:"id" bson:"id"`
	LastNotified time.Time `json:"lastNotified" bson:"lastNotified"`
}

type UserModel struct {
	DbClient *mongo.Client
}

func NewUserModel(client *mongo.Client) *UserModel {
	return &UserModel{
		DbClient: client,
	}
}

func (um *UserModel) UpdateSubscription(email string, subscription AlbumSubscription) error {

	coll := um.DbClient.Database("Notify").Collection("users")

	filter := bson.D{
		{
			Key: "email", Value: email,
		},
	}

	subscriptionExists := bson.M{
		"$elemMatch": bson.M{
			"id": subscription.Id,
		},
	}

	res := coll.FindOne(context.TODO(), bson.M{"email": email, "subscriptions": subscriptionExists})

	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		update := bson.M{"$push": bson.M{"subscriptions": subscription}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
	} else {
		update := bson.M{"$pull": bson.M{"subscriptions": bson.M{"id": subscription.Id}}}
		_, err := coll.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func (um *UserModel) RemoveSubscription(email string, albumId string) (string, error) {
	filter := bson.D{
		{
			Key: "email", Value: email,
		},
	}

	update := bson.M{"$pull": bson.M{"subscriptions": bson.M{"id": albumId}}}

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
		return User{Name: name, Email: email}, nil
	}
	var user User
	err := res.Decode(&user)
	if err != nil {
		fmt.Println("failed to decode user")
		fmt.Println(err.Error())
	}
	return user, nil
}

func (um *UserModel) ActivateSubscriptions(userName string, chatId int) error {

	filter := bson.M{
		"name": userName,
	}

	update := bson.M{"$set": bson.M{"chatId": chatId}}

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
