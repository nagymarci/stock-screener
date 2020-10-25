package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nagymarci/stock-screener/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Watchlists struct {
	collection *mongo.Collection
}

type WatchlistCollection interface {
	Create(watchlist model.WatchlistRequest) (primitive.ObjectID, error)
	Get(id primitive.ObjectID) (model.Watchlist, error)
	Delete(id primitive.ObjectID) (int64, error)
	GetAll(email string) ([]model.Watchlist, error)
	List() ([]model.Watchlist, error)
}

func NewWatchlists(db *mongo.Database) WatchlistCollection {
	return &Watchlists{
		collection: db.Collection("watchlist"),
	}
}

func (w *Watchlists) Create(watchlist model.WatchlistRequest) (primitive.ObjectID, error) {
	result, err := w.collection.InsertOne(context.TODO(), watchlist)

	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (w *Watchlists) Get(id primitive.ObjectID) (model.Watchlist, error) {
	var result model.Watchlist

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	err := w.collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func (w *Watchlists) Delete(id primitive.ObjectID) (int64, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	result, err := w.collection.DeleteOne(context.TODO(), filter)

	return result.DeletedCount, err
}

func (w *Watchlists) GetAll(email string) ([]model.Watchlist, error) {
	filter := bson.D{{Key: "email", Value: email}}

	cursor, err := w.collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var result []model.Watchlist

	for cursor.Next(context.TODO()) {
		var data model.Watchlist
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result, err
}

func (w *Watchlists) List() ([]model.Watchlist, error) {
	cursor, err := w.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	var result []model.Watchlist
	for cursor.Next(context.TODO()) {
		var data model.Watchlist
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result, err
}
