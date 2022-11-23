package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Mongo are the struct that allow interact with the mongo database.
type Mongo struct {
	collection *mongodriver.Collection
}

// NewMongo create mongo database connection.
func NewMongo(uri, database, collection, username, password string) (Mongo, error) {
	var m Mongo

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
		return Mongo{}, err
	}

	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		return Mongo{}, err
	}

	m.collection = cli.Database(database).Collection(collection)

	return m, nil
}

// Insert a new document, the doc param can be marshaled.
func (m Mongo) Insert(ctx context.Context, doc interface{}) error {
	_, err := m.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

// TODO: define this comment
// GetAll definition of this not already defined.
//
//nolint:godox,revive // are going be defined when the feature of get images are finished.
func (m Mongo) GetAll(ctx context.Context, query, projections bson.M) ([]Document, error) {
	opts := options.Find().SetProjection(projections)
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	c, err := m.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	res := make([]Document, c.RemainingBatchLength())

	for i := 0; c.Next(ctx); i++ {
		var doc Document
		if err := c.Decode(&doc); err != nil {
			return nil, err
		}

		res[i] = doc
	}

	return res, nil
}
