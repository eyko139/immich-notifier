package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/env"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"io"
	"net/http"
	"time"
)

type Immich struct {
	url string
}

type Album struct {
	AlbumName    string    `json:"albumName" bson:"albumName"`
	Description  string    `json:"description" bson:"description"`
	Id           string    `json:"id" bson:"id"`
	AlbumUsers   string    `json:"albumUsers" bson:"albumUsers"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
	LastNotified time.Time `json:"lastNotified" bson:"lastNotified"`
	IsSubscribed bool      `json:"isSubscribed" bson:"isSubscribed"`
}

type ImmichModel struct {
	DbClient *mongo.Client
	env      *env.Env
}

func NewImmichModel(client *mongo.Client, env *env.Env) *ImmichModel {
	return &ImmichModel{
		DbClient: client,
		env:      env,
	}
}

func (im *ImmichModel) FetchAlbums(apiKey string) ([]Album, error) {
	var albums []Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl+"/api/albums", nil)
	if err != nil {
		fmt.Println("Error creating request")
	}
	req.Header.Add("x-api-key", apiKey)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error executing request")
	}
	resBytes, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(resBytes, &albums)
	defer res.Body.Close()
	return albums, nil
}

func (im *ImmichModel) FetchAlbumsDetails(albumId, apiKey string) (Album, error) {
	var album Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl+"/api/albums/"+albumId, nil)
	if err != nil {
		fmt.Println("Error creating request")
	}
	req.Header.Add("x-api-key", apiKey)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error executing request")
	}
	resBytes, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(resBytes, &album)
	defer res.Body.Close()
	return album, nil
}

func (im *ImmichModel) InsertAlbum(album Album) {
	_, err := im.DbClient.Database("Notify").Collection("albums").InsertOne(context.TODO(), album, nil)
	if err != nil {
		fmt.Printf("Error saving album: %s", err)
	}
}

func (im *ImmichModel) UpdateSubscription(user User) {

	update := bson.D{{"$set", bson.D{{"Subscriptions", user.Subscriptions}}}}

	_, err := im.DbClient.Database("Notify").Collection("users").UpdateByID(context.TODO(), user.ID, update)
	if err != nil {
		fmt.Printf("Failed to update: %s", err)
	}
}
