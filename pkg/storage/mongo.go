package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoOptions struct {
	Host        string
	Port        string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
	MaxIdleTime time.Duration
	Schema      string
	Logger      *zap.Logger
}

func NewMongo(o MongoOptions) (*mongo.Database, error) {
	opts := options.Client()
	opts.SetMaxPoolSize(o.MaxPoolSize)
	opts.SetMinPoolSize(o.MinPoolSize)
	opts.SetTimeout(o.Timeout)
	opts.SetMaxConnIdleTime(o.MaxIdleTime)
	opts.ApplyURI("mongodb://" + o.Host + ":" + o.Port)
	client, err := mongo.NewClient(opts)
	if err != nil {
		o.Logger.Error("failed to settle mongo client", zap.Error(err))
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		o.Logger.Error("failed to settle mongo connection", zap.Error(err))
		panic(err)
	}
	return client.Database(o.Schema), nil
}
