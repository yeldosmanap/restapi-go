package repository

import (
	"context"
	"log"

	"gorest-api/internal/config"
	"gorest-api/internal/logs"
	"gorest-api/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectsDatabaseMongo struct {
	db *mongo.Collection
}

func NewProjectsDB(db *mongo.Database) *ProjectsDatabaseMongo {
	return &ProjectsDatabaseMongo{
		db: db.Collection("projects"),
	}
}

func (r *ProjectsDatabaseMongo) Create(ctx context.Context, project model.Project) (string, error) {
	logs.Log().Info("Creating a project in database...")

	result, err := r.db.InsertOne(ctx, project)

	if config.IsDuplicate(err) {
		logs.Log().Errorf("Duplicate project found %s", err.Error())
		id := primitive.NilObjectID.Hex()
		return id, model.ErrProjectIsAlreadyExists
	}

	log.Println(result.InsertedID.(primitive.ObjectID).Hex())

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return id, err
}

func (r *ProjectsDatabaseMongo) GetAll(ctx context.Context) ([]model.Project, error) {
	logs.Log().Info("Getting all projects from database...")

	query := bson.D{{}}

	cursor, err := r.db.Find(ctx, query)
	if err != nil {
		logs.Log().Errorf("Projects are not found in the database %s", err.Error())
		return []model.Project{}, model.ErrProjectNotFound
	}

	var projects = make([]model.Project, 0)

	logs.Log().Info("Extracting data from cursor in database...")

	if err := cursor.All(ctx, &projects); err != nil {
		logs.Log().Warn("Error occurred: %s", err.Error())
		return []model.Project{}, model.ErrProjectNotFound
	}

	return projects, nil
}

func (r *ProjectsDatabaseMongo) GetByTitle(ctx context.Context, userId string, title string) (model.Project, error) {
	logs.Log().Info("Getting all projects... ")

	var project model.Project

	filter := bson.M{"title": title, "user_id": userId}

	err := r.db.FindOne(ctx, filter).Decode(&project)

	return project, err
}

func (r *ProjectsDatabaseMongo) Delete(ctx context.Context, userId string, projectId string) error {
	logs.Log().Info("Deleting a project... ")

	projectObjectId, err := primitive.ObjectIDFromHex(projectId)
	if err != nil {
		return model.ErrCouldParseId
	}

	query := bson.M{"user_id": userId, "_id": projectObjectId}

	res, err := r.db.DeleteOne(ctx, &query)

	if res.DeletedCount < 1 {
		return model.ErrProjectNotFound
	}

	return err
}

func (r *ProjectsDatabaseMongo) Update(ctx context.Context, userId string, projectId string, newTitle string) error {
	updateQuery := bson.M{"title": newTitle, "user_id": userId}

	_, err := r.db.UpdateOne(ctx, bson.M{"_id": projectId, "user_id": userId}, bson.M{"$set": updateQuery})

	return err
}
