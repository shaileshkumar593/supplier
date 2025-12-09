package repository

import (
	"context"
	"fmt"
	"strings"

	svc "swallow-supplier/iface"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// mongoRepository implements the MongoRepository interface
type mongoRepository struct {
	client  *mongo.Client
	session mongo.Session
	db      *mongo.Database
	logger  log.Logger
}

func NewMongo(client *mongo.Client, databaseName string, logger log.Logger) (svc.MongoRepository, error) {
	//fmt.Println("*************** Database Name: ", databaseName)

	databaseName = strings.Trim(databaseName, "\"")

	if databaseName == "" || strings.Contains(databaseName, ".") {
		return nil, fmt.Errorf("invalid database name: %s", databaseName)
	}
	db := client.Database(databaseName)

	// Check if the database name is correct
	if databaseName == "" {
		return nil, fmt.Errorf("invalid database name: database name cannot be empty")
	}

	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	return &mongoRepository{
		client:  client,
		session: session,
		db:      db,
		logger:  log.With(logger, "repo", "mongo"),
	}, nil
}

/* // GetMongoClient returns a singleton instance of the MongoDB client
func GetMongoClient(mongoURI string) (*mongo.Client, error) {
	var err error
	mongoClientOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(mongoURI)
		mongoClientInstance, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			err = fmt.Errorf("failed to connect to MongoDB: %w", err)
		}
	})
	return mongoClientInstance, err
} */

func (r *mongoRepository) GetDbTx(ctx context.Context) (mongo.Session, error) {
	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *mongoRepository) CommitTransaction(ctx context.Context) error {
	return r.session.CommitTransaction(ctx)
}

func (r *mongoRepository) AbortTransaction(ctx context.Context) error {
	return r.session.AbortTransaction(ctx)
}

func (r *mongoRepository) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := r.db.Collection(collection)
	return coll.InsertOne(ctx, document)
}

func (r *mongoRepository) InsertMany(ctx context.Context, collection string, documents []interface{}) (*mongo.InsertManyResult, error) {
	coll := r.db.Collection(collection)
	return coll.InsertMany(ctx, documents)
}

func (r *mongoRepository) Find(ctx context.Context, collection string, filter interface{}) (bson.M, error) {
	coll := r.db.Collection(collection)
	var result bson.M
	err := coll.FindOne(ctx, filter).Decode(&result)
	return result, err
}

func (r *mongoRepository) FindMany(ctx context.Context, collection string, filter interface{}) ([]bson.M, error) {
	coll := r.db.Collection(collection)
	var results []bson.M
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *mongoRepository) Update(ctx context.Context, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := r.db.Collection(collection)
	return coll.UpdateOne(ctx, filter, update)
}

func (r *mongoRepository) UpdateMany(ctx context.Context, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	coll := r.db.Collection(collection)
	return coll.UpdateMany(ctx, filter, update)
}

func (r *mongoRepository) Delete(ctx context.Context, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	coll := r.db.Collection(collection)
	return coll.DeleteOne(ctx, filter)
}

// GetMongoClient returns the underlying MongoDB client.
func (r *mongoRepository) GetMongoClient(ctx context.Context) *mongo.Client {
	return r.client
}

func (r *mongoRepository) GetMongoDb(ctx context.Context) *mongo.Database {
	return r.db
}
