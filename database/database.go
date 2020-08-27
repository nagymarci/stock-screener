package database

import (
	"context"
	"log"
	"time"

	"github.com/nagymarci/stock-screener/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var database *mongo.Database

//Connect establishes the connection to the database
func Connect(connectionURI string) {
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())

	if err != nil {
		log.Fatal("Couldn't connect to database", err)
	} else {
		log.Println("Connected to database")
	}

	database = client.Database("stock-screener")
}

//Save writes the stockData to the database
func Save(stockData model.StockDataInfo) {
	collection := database.Collection("stockinfo")

	insertedResult, err := collection.InsertOne(context.TODO(), stockData)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("stockData inserted into DB", insertedResult)

}

//Get retreives the stockinfo for the given symbol
func Get(symbol string) model.StockDataInfo {
	collection := database.Collection("stockinfo")

	var result model.StockDataInfo

	filter := bson.D{primitive.E{Key: "ticker", Value: symbol}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		log.Println(err)
	}

	return result
}

//GetAll retreives all of the objects from the database
func GetAll() []model.StockDataInfo {
	collection := database.Collection("stockinfo")

	var result []model.StockDataInfo

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	for cursor.Next(context.TODO()) {
		var data model.StockDataInfo
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result
}

//GetAllExpired returns list of stocks that has at least one value expired
func GetAllExpired() []model.StockDataInfo {
	collection := database.Collection("stockinfo")

	var result []model.StockDataInfo

	now := time.Now()

	filter := bson.D{{"$or", bson.A{
		bson.D{{"nextUpdate", bson.D{{"$lt", now}}}},
		bson.D{{"dividendYield5yr.nextUpdate", bson.D{{"$lt", now}}}},
		bson.D{{"peRatio5yr.nextUpdate", bson.D{{"$lt", now}}}},
		bson.D{{"nextUpdate", nil}}}}}

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		log.Fatalln(err)
	}

	for cursor.Next(context.TODO()) {
		var data model.StockDataInfo
		cursor.Decode(&data)
		result = append(result, data)
	}

	return result
}
