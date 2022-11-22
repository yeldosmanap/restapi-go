package config

import (
	"context"
	"errors"

	"github.com/spf13/viper"

	"gorestapi/internal/apperror"
	"gorestapi/internal/logs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	URI      string `mapstructure:"mongoURI"`
	Username string `mapstructure:"mongoUsername"`
	Password string `mapstructure:"mongoPassword"`
	Name     string `mapstructure:"databaseName"`
}

func MongoNewClient(ctx context.Context, cancel context.CancelFunc, mongoCfg *MongoConfig) (*mongo.Client, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	err := viper.UnmarshalKey("mongo", &mongoCfg)
	if err != nil {
		return nil, err
	}

	opts := options.Client().ApplyURI(mongoCfg.URI)

	if mongoCfg.Username != "" && mongoCfg.Password != "" {
		opts.SetAuth(
			options.Credential{
				Username: mongoCfg.Username,
				Password: mongoCfg.Password,
			})
	} else {
		return nil, apperror.ErrBadCredentials
	}

	logs.Log().Info("Enabling new mongodb client")
	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, err
	}

	defer cancel()

	logs.Log().Info("Connecting to the database")
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	logs.Log().Info("Pinging the database")
	err = client.Ping(ctx, nil)
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
