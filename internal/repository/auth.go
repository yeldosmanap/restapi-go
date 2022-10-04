package repository

import (
	"context"

	"gorest-api/internal/logs"
	"gorest-api/internal/model"
	"gorest-api/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthPostgres struct {
	db *mongo.Collection
}

func NewAuthMongo(db *mongo.Database) *AuthPostgres {
	return &AuthPostgres{
		db: db.Collection("users"),
	}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (string, error) {
	logs.Log().Info("Creating a user in database...")
	result, err := r.db.InsertOne(ctx, user)

	id := result.InsertedID.(primitive.ObjectID).Hex()

	logs.Log().Info("Inserted a user with ID: %s", id)
	if err != nil {
		logs.Log().Info("User with this id already exists: %s", err.Error())
		return "", utils.ErrEmailAlreadyExists
	}

	return id, err
}

func (r *AuthPostgres) GetUser(ctx context.Context, email, password string) (model.User, error) {
	logs.Log().Info("Getting a user from database...")

	var user model.User

	if err := r.db.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		logs.Log().Info("Error occurred: %s", err.Error())
		return model.User{}, model.ErrUserNotFound
	}

	return user, nil
}

func (r *AuthPostgres) GetById(ctx context.Context, id string) (model.User, error) {
	userId, err := primitive.ObjectIDFromHex(id)
	var user model.User

	err = r.db.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)

	return user, err
}

func (r *AuthPostgres) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var userFound model.User

	err := r.db.FindOne(ctx, bson.M{"email": email}).Decode(&userFound)

	return userFound, err
}

func (r *AuthPostgres) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User

	result, err := r.db.Find(ctx, bson.M{})

	for result.Next(ctx) {
		var user model.User

		if err = result.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}
