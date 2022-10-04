package config

import (
	"context"
	"errors"
	"time"

	"gorest-api/internal/logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10 * time.Second

func NewClient(uri, username, password string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)

	if username != "" && password != "" {
		opts.SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			})
	}

	logs.Log().Info("Enabling new client")
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	logs.Log().Info("Enabling context")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logs.Log().Info("Connecting to the database")
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	logs.Log().Info("Pinging the database")
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}
