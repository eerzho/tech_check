package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(db, url string) (*mongo.Database, error) {
	const op = "util.NewMongo"

	clientOPT := options.Client().ApplyURI(url).SetMaxPoolSize(1)

	var err error
	var client *mongo.Client

	attempts := 10
	for attempts > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, clientOPT)
		if err == nil {
			err = client.Ping(ctx, nil)
			if err == nil {
				break
			}
		}
		log.Printf("%s: attempts left - %d", op, attempts)
		time.Sleep(time.Second)
		attempts--
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return client.Database(db), nil
}
