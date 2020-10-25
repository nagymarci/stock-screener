package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type recommendations struct{}

var recommendationCollection *mongo.Collection

//Recommendations gives CRUD operations for recommendations
var Recommendations recommendations

type recommendation struct {
	ID     primitive.ObjectID `bson:"_id"`
	Stocks []string           `bson:"stocks"`
}

func (recommendations) Create(id primitive.ObjectID, stocks []string) error {
	_, err := watchlistCollection.InsertOne(context.TODO(), recommendation{ID: id, Stocks: stocks})

	return err
}

func (recommendations) Get(id primitive.ObjectID) ([]string, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var result recommendation
	err := watchlistCollection.FindOne(context.TODO(), filter).Decode(&result)

	return result.Stocks, err
}

func (recommendations) Update(id primitive.ObjectID, stocks []string) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	opts := options.Replace().SetUpsert(true)

	_, err := watchlistCollection.ReplaceOne(context.TODO(), filter, recommendation{ID: id, Stocks: stocks}, opts)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("recommendation inserted into DB ", id)
	return nil
}
