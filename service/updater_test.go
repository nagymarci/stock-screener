package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nagymarci/stock-screener/model"
	"github.com/nagymarci/stock-screener/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/nagymarci/stock-screener/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func TestMain(m *testing.M) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(time.Minute * 2),
		Env:          map[string]string{},
	}
	req.Env["MONGO_INITDB_ROOT_USERNAME"] = "mongodb"
	req.Env["MONGO_INITDB_ROOT_PASSWORD"] = "mongodb"
	req.Env["MONGO_INITDB_DATABASE"] = "stock-screener"

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatalln(err)
	}
	defer mongoC.Terminate(ctx)
	ip, err := mongoC.Host(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		log.Fatalln(err)
	}

	dbConnectionURI := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d",
		"mongodb",
		"mongodb",
		ip,
		port.Int())

	db = database.New(dbConnectionURI)

	code := m.Run()

	os.Exit(code)
}

func TestUpdater(t *testing.T) {
	t.Run("updates stock when nextUpdate is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		stockData := model.StockDataInfo{}
		stockData.Ticker = "INTC"
		stockData.Dividend = 0.33
		stockData.Eps = 5.43
		stockData.Price = 49.28
		stockData.DividendYield5yr.Avg = 2.62
		stockData.DividendYield5yr.Max = 3.65
		stockData.PeRatio5yr.Avg = 14.89
		stockData.PeRatio5yr.Min = 8.79

		sDb := database.NewStockinfos(db)

		err := sDb.Save(stockData)
		if err != nil {
			t.Fatal(err)
		}
		defer sDb.Delete(stockData.Ticker)

		sSC := mocks.NewMockgetStockWithFields(ctrl)
		stockData.Price = 100
		sSC.EXPECT().GetWithFields("INTC", []string{"price", "eps", "div", "divHist", "pe"}).Return(stockData, nil)

		updater := New(sDb, sSC, "1h", "1h", "1h")

		updater.UpdateStocks()

		result, err := sDb.Get(stockData.Ticker)

		if err != nil {
			t.Fatal(err)
		}

		if result.Price != 100 {
			t.Fatalf("stock is not updated")
		}
	})
	t.Run("updates pe when pe.nextUpdate is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		stockData := model.StockDataInfo{}
		stockData.Ticker = "INTC"
		stockData.Dividend = 0.33
		stockData.Eps = 5.43
		stockData.Price = 49.28
		stockData.DividendYield5yr.Avg = 2.62
		stockData.DividendYield5yr.Max = 3.65
		stockData.PeRatio5yr.Avg = 14.89
		stockData.PeRatio5yr.Min = 8.79
		stockData.NextUpdate = time.Now().Add(5000000000)
		stockData.DividendYield5yr.NextUpdate = time.Now().Add(5000000000)

		sDb := database.NewStockinfos(db)

		err := sDb.Save(stockData)
		if err != nil {
			t.Fatal(err)
		}
		defer sDb.Delete(stockData.Ticker)

		sSC := mocks.NewMockgetStockWithFields(ctrl)
		stockData.Price = 100
		sSC.EXPECT().GetWithFields("INTC", []string{"pe"}).Return(stockData, nil)

		updater := New(sDb, sSC, "1h", "1h", "1h")

		updater.UpdateStocks()

		result, err := sDb.Get(stockData.Ticker)

		if err != nil {
			t.Fatal(err)
		}

		if result.Price != 100 {
			t.Fatalf("stock is not updated")
		}
	})
}
