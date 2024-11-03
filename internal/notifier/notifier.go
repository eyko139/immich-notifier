package notifier

import (
	"context"
	"fmt"
	"github.com/eyko139/immich-notifier/internal/env"
	"github.com/eyko139/immich-notifier/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type Notifier struct {
	interval time.Duration
	client   *mongo.Client
	env      *env.Env
	immich   *models.ImmichModel
}

func New(client *mongo.Client, env *env.Env, interval time.Duration, immich *models.ImmichModel) *Notifier {
	return &Notifier{
		interval: interval,
		client:   client,
		env:      env,
		immich:   immich,
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
		for _, val := range result {
			for _, albumId := range val.SubscribedAlbums {
				if len(albumId) == 0 {
					continue
				}
				album, err := n.immich.FetchAlbumsDetails(albumId, val.ApiKey)
				if err != nil {
					fmt.Println("error fetching album")
				}
				fmt.Println(album)
			}
		}
	}
}
