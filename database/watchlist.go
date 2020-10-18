package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nagymarci/stock-screener/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type watchlists struct{}

var watchlistCollection *mongo.Collection

//WatchLists gives CRUD operations for watchlists
var WatchLists watchlists

func (watchlists) Create(watchlist model.WatchlistRequest) (primitive.ObjectID, error) {
	result, err := watchlistCollection.InsertOne(context.TODO(), watchlist)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (watchlists) Get(id primitive.ObjectID) (model.Watchlist, error) {
	var result model.Watchlist

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	err := watchlistCollection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (watchlists) Delete(id primitive.ObjectID) (int64, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	result, err := watchlistCollection.DeleteOne(context.TODO(), filter)

	return result.DeletedCount, err
}
