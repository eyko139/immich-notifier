package models

import (
	"encoding/json"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/env"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"io"
	"net/http"
	"time"
)

type Immich struct {
	url string
}

type Album struct {
	AlbumName   string `json:"albumName"`
	Description string `json:"description"`
	Id          string `json:"id"`
	AlbumUsers  string `json:"albumUsers"`
}

type ImmichModel struct {
	DbClient *mongo.Client
	env      *env.Env
}

func NewImmichModel(client *mongo.Client, env *env.Env) *ImmichModel {
	return &ImmichModel{
		DbClient: client,
		env: env,
	}
}

func (im *ImmichModel) FetchAlbums(apiKey string) ([]Album, error) {
	var albums []Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl + "/api/albums", nil)
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
	fmt.Println(albums)
	return albums, nil
}

func (im *ImmichModel) FetchAlbumsDetails(albumId, apiKey string) (Album, error) {
	var album Album
	req, err := http.NewRequest(http.MethodGet, im.env.ImmichUrl + "/api/album/" + albumId, nil)
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
	fmt.Println(album)
	return album, nil
}
