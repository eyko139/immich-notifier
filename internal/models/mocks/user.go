package mocks

import (
	"fmt"

	"github.com/eyko139/immich-notifier/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserModel struct{}

var mockSubscription = models.AlbumSubscription{}

func newMockUser() *models.User {
	mockId, _ := bson.ObjectIDFromHex("123")

	var mockUser = &models.User{
		Name:          "mockUser",
		Subscriptions: []models.AlbumSubscription{mockSubscription},
		Email:         "test@test.de",
		ID:            mockId,
		ApiKey:        "123",
		ChatId:        0,
	}

	return mockUser
}

func (um *UserModel) UpdateSubscription(email string, subscription models.AlbumSubscription) error {
	return nil
}

func (um *UserModel) FindOrInsertUser(name, email string) (models.User, error) {
	mockuser := newMockUser()
    fmt.Println(name)
    if name == "active" {
        mockuser.Subscriptions[0].Id = "1"
    }
	return *mockuser, nil
}

func (um *UserModel) RemoveSubscription(email, album string) (string, error) {
	return "", nil
}

func (um *UserModel) ActivateSubscriptions(userId string, chatId int) error {
	return nil
}
