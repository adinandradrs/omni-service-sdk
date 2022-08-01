package storage

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongo struct {
	Host        string
	Port        string
	MaxPoolSize uint64
	MinPoolSize uint64
	Timeout     time.Duration
	MaxIdleTime time.Duration
	Schema      string
	Logger      *zap.Logger
}

func NewMongoClient(m Mongo) (*mongo.Database, error) {
	opts := options.Client()
	opts.SetMaxPoolSize(m.MaxPoolSize)
	opts.SetMinPoolSize(m.MinPoolSize)
	opts.SetTimeout(m.Timeout)
	opts.SetMaxConnIdleTime(m.MaxIdleTime)
	opts.ApplyURI("mongodb://" + m.Host + ":" + m.Port)
	client, err := mongo.NewClient(opts)
	if err != nil {
		m.Logger.Error("failed to settle mongo client", zap.Error(err))
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		m.Logger.Error("failed to settle mongo connection", zap.Error(err))
		panic(err)
	}
	return client.Database(m.Schema), nil
}
