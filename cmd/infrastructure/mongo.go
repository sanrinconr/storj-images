package infrastructure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo is the interface to interact with the ids database.
type Mongo interface {
	Insert(context.Context, interface{}) error
	GetAll(context.Context, bson.M, bson.M) ([]bson.M, error)
}

type mongo struct {
	collection *mongodriver.Collection
}

// NewMongo create mongo database connection.
func NewMongo(uri, database, collection, username, password string) (Mongo, error) {
	var m mongo

	cli, err := mongodriver.Connect(
		context.TODO(),
		options.Client().ApplyURI(uri).SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	m.collection = cli.Database(database).Collection(collection)

	return m, nil
}

func (m mongo) Insert(ctx context.Context, doc interface{}) error {
	_, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return err
}

func (m mongo) GetAll(ctx context.Context, query, projections bson.M) ([]bson.M, error) {
	opts := options.Find().SetProjection(projections)

	c, err := m.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := c.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, err
}
