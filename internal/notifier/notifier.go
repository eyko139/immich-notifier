package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log"
	"net/http"
	"time"
)

type Notifier struct {
	interval time.Duration
	client   *mongo.Client
	env      *env.Env
	immich   *models.ImmichModel
	errLog   *log.Logger
	infoLog  *log.Logger
}

type Notification struct {
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority int    `json:"priority"`
}

func New(client *mongo.Client, env *env.Env, interval time.Duration, immich *models.ImmichModel, errLog *log.Logger, infoLog *log.Logger) *Notifier {
	return &Notifier{
		interval: interval,
		client:   client,
		env:      env,
		immich:   immich,
		errLog:   errLog,
		infoLog:  infoLog,
	}
}

func (n *Notifier) StartLoop() {

	ticker := time.NewTicker(n.interval)
	var result []models.User

	for {
		<-ticker.C
		fmt.Println("Ticker ticket")
		cursor, err := n.client.Database("Notify").Collection("users").Find(context.TODO(), bson.D{}, nil)
		if err != nil {
			fmt.Println(err)
		}
		if err := cursor.All(context.TODO(), &result); err != nil {
			fmt.Println("Error unpacking cursor")
		}
		for _, user := range result {
			for idx, subscription := range user.Subscriptions {
				if !subscription.IsSubscribed {
					continue
				}
				if err != nil {
					fmt.Println("error fetching album")
				}
				album, _ := n.immich.FetchAlbumsDetails(subscription.Id)
				fmt.Printf("checking dates: albumUpdate: %s, subscriptionLastNotified: %s", album.UpdatedAt, subscription.LastNotified)
				if album.UpdatedAt.After(subscription.LastNotified) {
					user.Subscriptions[idx].LastNotified = time.Now()
					n.immich.UpdateSubscription(user)
					n.Notify(user, subscription, album.Assets[0].ID)
				}
			}
		}
	}
}

func (n *Notifier) Notify(user models.User, sub models.AlbumSubscription, latestAssetId string) {
	thumbBytes := n.immich.FetchThumbnail(latestAssetId)
	n.Gotify(user, sub)
	n.Telegram(user, thumbBytes, sub)
}

func (n *Notifier) Telegram(user models.User, latestAssetBytes []byte, album models.AlbumSubscription) {

	thumbNailRequest := buildThumbnailRequest(latestAssetBytes, user.ChatId, album)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	thumbResponse, err := client.Do(thumbNailRequest)

	if err != nil {
		n.errLog.Println("Error sending thumbnail" + err.Error())
	}

	n.infoLog.Printf("Sent update thumbnail %+v", thumbResponse)
}

func (n *Notifier) Gotify(user models.User, sub models.AlbumSubscription) {
	notification := Notification{
		Message:  fmt.Sprintf("Album %s has been updated, user: %s", sub.AlbumName, user.Email),
		Title:    "Immich album update",
		Priority: 1,
	}

	notificationBytes, _ := json.Marshal(notification)

	req, _ := http.NewRequest(http.MethodPost, n.env.GotifyUrl, bytes.NewBuffer(notificationBytes))
	req.Header.Set(GotifyAuthHeader, n.env.GotifyKey)
	req.Header.Set(ContentType, JsonContentType)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to notify: %s", err)
	}
	n.infoLog.Printf("Sent gotify notification, res: %v", res)
}
