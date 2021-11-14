package database

import (
	"context"
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MongoConnect() (*mongo.Database, error) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/TradingBlack")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb")

	return client.Database("user"), nil
}

func GormCon() *gorm.DB {
	host := "localhost"
	user := "postgres"
	password := "root"
	dbname := "other"
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=EST", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("there was an error connecting through gorm")
		log.Println(err.Error())
	}
	return db
}

func PGCon() (con *pg.DB) {
	address := fmt.Sprintf("%s:%s", "localhost", "5432")
	options := &pg.Options{
		User:     "postgres",
		Password: "root",
		Addr:     address,
		Database: "other",
	}

	con = pg.Connect(options)
	if con == nil {
		log.Fatal("could not connect to pg")
	}
	return con
}
