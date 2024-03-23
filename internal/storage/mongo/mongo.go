package mongo

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go_db/config"
	"go_db/internal/models"
)

type MongoStore struct {
	log  *slog.Logger
	cl   *mongo.Client
	coll *mongo.Collection
}

func GetStore(ctx context.Context, logger *slog.Logger) (*MongoStore, error) {
	logger.Info("connecting DB mongo...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetMongoString()))
	if err != nil {
		return nil, err
	}

	i := MongoStore{log: logger, cl: client}
	i.log.Info("DB mongo connected")
	i.coll = i.cl.Database("Empl").Collection("Employee")

	return &i, nil
}

func (i MongoStore) GetConnection(ctx context.Context) (*sql.DB, error) {
	return nil, errors.New("mongo: GetConnection not supported")
}

func (i MongoStore) Release(ctx context.Context) {

	err := i.cl.Disconnect(context.TODO())
	if err != nil {
		i.log.Error("DB mongo disconnect error:", slog.String("message", err.Error()))
	} else {
		i.log.Info("DB mongo disconnected")
	}
}

func (i MongoStore) EmployeeCreate(ctx context.Context, empl models.Employee) (uint32, error) {
	i.log.Debug("mongo:create ", slog.String("Name", empl.Name))
	id64, err := i.coll.CountDocuments(ctx, bson.D{{}})
	if err != nil {
		return 0, err
	}
	empl.Id = uint32(id64) + 1

	result, err := i.coll.InsertOne(ctx, empl)
	if err != nil {
		return 0, err
	}
	i.log.Debug("Inserteds", slog.Any("Id", result.InsertedID))
	return empl.Id, nil
}

func (i MongoStore) EmployeeGet(ctx context.Context, id uint32) (*models.Employee, error) {
	i.log.Debug("mongo:get ", slog.Any("ID", id))
	var empl models.Employee
	err := i.coll.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&empl)
	return &empl, err
}
func (i MongoStore) EmployeeUpdate(ctx context.Context, empl models.Employee) error {
	i.log.Debug("mongo:update ", slog.Any("ID", empl.Id), slog.String("Name", empl.Name))

	filter := bson.D{{Key: "id", Value: empl.Id}}
	val := bson.D{{Key: "$set", Value: bson.D{{Key: "id", Value: empl.Id}, {Key: "Name", Value: empl.Name}}}}

	updateResult, err := i.coll.UpdateOne(ctx, filter, val)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount != 1 {
		return errors.New("not found")
	}
	if updateResult.ModifiedCount != 1 {
		return errors.New(" found, but no updated")
	}
	return nil
}
func (i MongoStore) EmployeeDelete(ctx context.Context, id uint32) error {
	i.log.Debug("mongo:delete ", slog.Any("ID", id))
	filter := bson.D{{Key: "id", Value: id}}
	deleteResult, err := i.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount != 1 {
		return errors.New("not found")
	}
	return nil
}
func (i MongoStore) EmployeeList(ctx context.Context) ([]*models.Employee, error) {
	i.log.Debug("mongo:list")
	cursor, err := i.coll.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	var result []*models.Employee

	//err = cursor.Decode(&rx)
	for cursor.Next(ctx) {
		var elem models.Employee
		err := cursor.Decode(&elem)
		if err != nil {
			return nil, err
		}

		result = append(result, &elem)
	}
	return result, nil
}
