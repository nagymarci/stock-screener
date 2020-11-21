package database

import (
	"context"
	"time"

	"github.com/nagymarci/stock-screener/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Stockinfos struct {
	collection *mongo.Collection
}

func NewStockinfos(db *mongo.Database) *Stockinfos {
	return &Stockinfos{
		collection: db.Collection("stockinfo"),
	}
}

//Save writes the stockData to the database
func (si *Stockinfos) Save(stockData model.StockDataInfo) error {
	_, err := si.collection.InsertOne(context.TODO(), stockData)

	return err
}

//Update sets the fields that were changed in the DB
func (si *Stockinfos) Update(stockData model.StockDataInfo) error {
	filter := bson.D{{Key: "ticker", Value: stockData.Ticker}}

	update := bson.A{bson.D{{Key: "$set", Value: composeSetFields(&stockData)}}}

	_, err := si.collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func composeSetFields(stockData *model.StockDataInfo) bson.D {
	var setFields bson.D

	if stockData.Price != 0 || stockData.Eps != 0 || stockData.Dividend != 0 {
		setFields = append(setFields, bson.E{Key: "nextUpdate", Value: stockData.NextUpdate})
	}

	if stockData.Price != 0 {
		setFields = append(setFields, bson.E{Key: "price", Value: stockData.Price})
	}

	if stockData.Eps != 0 {
		setFields = append(setFields, bson.E{Key: "eps", Value: stockData.Eps})
	}

	if stockData.Dividend != 0 {
		setFields = append(setFields, bson.E{Key: "dividend", Value: stockData.Dividend})
	}

	if stockData.DividendYield5yr.Avg != 0 || stockData.DividendYield5yr.Max != 0 {
		setFields = append(setFields, bson.E{Key: "dividendYield5yr", Value: stockData.DividendYield5yr})
	}

	if stockData.PeRatio5yr.Avg != 0 || stockData.PeRatio5yr.Min != 0 {
		setFields = append(setFields, bson.E{Key: "peRatio5yr", Value: stockData.PeRatio5yr})
	}

	return setFields
}

//Get retreives the stockinfo for the given symbol
func (si *Stockinfos) Get(symbol string) (model.StockDataInfo, error) {
	var result model.StockDataInfo

	filter := bson.D{primitive.E{Key: "ticker", Value: symbol}}

	err := si.collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

//GetAll retreives all of the objects from the database
func (si *Stockinfos) GetAll() ([]model.StockDataInfo, error) {
	cursor, err := si.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	var result []model.StockDataInfo

	for cursor.Next(context.TODO()) {
		var data model.StockDataInfo
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result, err
}

//GetAllExpired returns list of stocks that has at least one value expired
func (si *Stockinfos) GetAllExpired() ([]model.StockDataInfo, error) {
	now := time.Now()

	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "dividendYield5yr.nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "peRatio5yr.nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "nextUpdate", Value: nil}}}}}

	cursor, err := si.collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var result []model.StockDataInfo

	for cursor.Next(context.TODO()) {
		var data model.StockDataInfo
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result, err
}

//Delete removes the given symbol from the database
func (si *Stockinfos) Delete(symbol string) error {
	filter := bson.D{{Key: "ticker", Value: symbol}}

	_, err := si.collection.DeleteOne(context.TODO(), filter)

	return err
}
