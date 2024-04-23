package repositories

import (
	"context"
	"github.com/goravel/framework/facades"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CollectionNameGetter interface {
	CollectionName() string
	DatabaseName() string
}

func Connection(c CollectionNameGetter) *mongo.Collection {
	clientOptions := options.Client().ApplyURI(facades.Config().GetString("DB_STRING", "")).SetTimeout(10 * time.Second)
	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		facades.Log().Error(err)
		return nil
	}
	return client.Database(c.DatabaseName()).Collection(c.CollectionName())
}
