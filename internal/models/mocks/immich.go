package mocks

import (
	"time"

	"github.com/eyko139/immich-notifier/internal/models"
)

type ImmichModel struct{}

func newMockAlbum() models.Album {

	var albumOwner = models.AlbumUser{
		Id:    "45",
		Email: "test1@test.de",
		Name:  "TestAlbumOwner",
	}
	var albumUser = models.AlbumUser{
		Id:    "44",
		Email: "test@test.de",
		Name:  "TestAlbumUser",
	}
	var assets = []struct {
		ID string `json:"id"`
	}{
		{ID: "15"},
	}

	var albumUsers = []struct {
		User models.AlbumUser `json:"user" bson:"user"`
	}{
		{User: albumUser},
	}

	return models.Album{
		AlbumName:             "mockAlbum",
		Description:           "very cool desc",
		Id:                    "1",
		UpdatedAt:             time.Now(),
		LastNotified:          time.Now(),
		IsSubscribed:          false,
		AlbumThumbnailAssetId: "1234",
		B64Thumbnail:          "thumbNailString",
		AssetCount:            17,
		Assets:                assets,
		Owner:                 albumOwner,
		AlbumUsers:            albumUsers,
	}
}

func (im *ImmichModel) FetchAlbums(userEmail string) ([]models.Album, error) {
	return []models.Album{newMockAlbum()}, nil
}
func (im *ImmichModel) FetchAlbumsDetails(albumId string) (*models.Album, error) {
	return nil, nil
}

func (im *ImmichModel) InsertOrAlbum(album models.Album) {
}

func (im *ImmichModel) UpdateSubscription(user models.User) {
}

func (im *ImmichModel) FetchThumbnail(uuid string) []byte {
	return []byte("base64mock bytes")
}
