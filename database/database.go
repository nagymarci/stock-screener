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
func Save(stockData model.StockDataInfo) error {
	collection := database.Collection("stockinfo")

	_, err := collection.InsertOne(context.TODO(), stockData)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("stockData inserted into DB ", stockData.Ticker)
	return nil
}

//Update sets the fields that were changed in the DB
func Update(stockData model.StockDataInfo) {
	collection := database.Collection("stockinfo")

	filter := bson.D{{Key: "ticker", Value: stockData.Ticker}}

	update := bson.A{bson.D{{Key: "$set", Value: composeSetFields(&stockData)}}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("stockData updated ", stockData.Ticker)
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
func Get(symbol string) (model.StockDataInfo, error) {
	collection := database.Collection("stockinfo")

	var result model.StockDataInfo

	filter := bson.D{primitive.E{Key: "ticker", Value: symbol}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
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

	filter := bson.D{{Key: "$or", Value: bson.A{
		bson.D{{Key: "nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "dividendYield5yr.nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "peRatio5yr.nextUpdate", Value: bson.D{{Key: "$lt", Value: now}}}},
		bson.D{{Key: "nextUpdate", Value: nil}}}}}

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

//Delete removes the given symbol from the database
func Delete(symbol string) error {
	collection := database.Collection("stockinfo")

	filter := bson.D{{Key: "ticker", Value: symbol}}

	_, err := collection.DeleteOne(context.TODO(), filter)

	return err
}

func DeleteProfile(profile string) error {
	collection := database.Collection("profiles")

	filter := bson.D{{Key: "name", Value: profile}}

	_, err := collection.DeleteOne(context.TODO(), filter)

	return err
}

func SaveProfile(profile model.Profile) error {
	collection := database.Collection("profiles")

	_, err := collection.InsertOne(context.TODO(), profile)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("profile inserted into DB ", profile.Name)
	return nil
}

func GetProfile(profile string) (model.Profile, error) {
	collection := database.Collection("profiles")

	var result model.Profile

	filter := bson.D{primitive.E{Key: "name", Value: profile}}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result, err
}

func GetAllProfileName() ([]string, error) {
	collection := database.Collection("profiles")

	var result []string

	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var data model.Profile
		cursor.Decode(&data)
		result = append(result, data.Name)
	}

	return result, err
}

func GetPreviouslyRecommendedStocks(profileName string) ([]string, error) {
	collection := database.Collection("recommendations")

	filter := bson.D{primitive.E{Key: "name", Value: profileName}}

	var result model.Profile

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	return result.Stocks, err
}

func SaveRecommendation(profileName string, stocks []string) error {
	collection := database.Collection("recommendations")

	filter := bson.D{primitive.E{Key: "name", Value: profileName}}
	opts := options.Replace().SetUpsert(true)

	_, err := collection.ReplaceOne(context.TODO(), filter, model.Profile{Name: profileName, Stocks: stocks}, opts)

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("recommendation inserted into DB ", profileName)
	return nil
}

func DeleteRecommendation(profileName string) error {
	collection := database.Collection("recommendations")

	filter := bson.D{{Key: "name", Value: profileName}}

	_, err := collection.DeleteOne(context.TODO(), filter)

	return err
}
