package models

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/util"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"io"
	"net/http"
	"time"
)

const (
	ImmichApiHeader = "x-api-key"
)

type Immich struct {
	url string
}

type Album struct {
	AlbumName             string    `json:"albumName" bson:"albumName"`
	Description           string    `json:"description" bson:"description"`
	Id                    string    `json:"id" bson:"id"`
	AlbumUsers            string    `json:"albumUsers" bson:"albumUsers"`
	UpdatedAt             time.Time `json:"updatedAt" bson:"updatedAt"`
	LastNotified          time.Time `json:"lastNotified" bson:"lastNotified"`
	IsSubscribed          bool      `json:"isSubscribed" bson:"isSubscribed"`
	AlbumThumbnailAssetId string    `json:"albumThumbnailAssetId" bson:"albumThumbnailAssetId"`
	B64Thumbnail          string    `json:"b64Thumbnail" bson:"b64Thumbnail"`
	AssetCount            int       `json:"assetCount" bson:"assetCount"`
	Assets                []struct {
		ID string `json:"id"`
	}
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

func (im *ImmichModel) FetchAlbums() ([]Album, error) {
	var albums []Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl+"/api/albums", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(ImmichApiHeader, im.env.ImmichApiKey)
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

	filteredAlbums := util.Filter(albums, IsNotEmpty)

	for idx, album := range filteredAlbums {
		thumbNail := im.FetchThumbnail(album.AlbumThumbnailAssetId)
		base64String := base64.StdEncoding.EncodeToString(thumbNail)
		filteredAlbums[idx].B64Thumbnail = base64String
	}

	return filteredAlbums, nil
}

func (im *ImmichModel) FetchAlbumsDetails(albumId string) (Album, error) {
	var album Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl+"/api/albums/"+albumId, nil)
	if err != nil {
		fmt.Println("Error creating request")
	}
	req.Header.Add(ImmichApiHeader, im.env.ImmichApiKey)
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

func (im *ImmichModel) InsertOrAlbum(album Album) {

	filter := bson.D{
		{
			"id", album.Id,
		},
	}
	res := im.DbClient.Database("Notify").Collection("albums").FindOneAndReplace(context.TODO(), filter, album, nil)

	if res.Err() == nil {
		return
	}

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

func (im *ImmichModel) FetchThumbnail(uuid string) []byte {

	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl+"/api/assets/"+uuid+"/thumbnail", nil)
	if err != nil {
		fmt.Println("Error creating request")
	}
	req.Header.Add(ImmichApiHeader, im.env.ImmichApiKey)
	req.Header.Add("Accept", "application/octet-stream")
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := client.Do(req)

	bytes, _ := io.ReadAll(res.Body)
	return bytes
}

func IsNotEmpty(album Album) bool {
	return album.AssetCount != 0
}
