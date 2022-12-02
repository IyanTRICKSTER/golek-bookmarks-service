package database

import (
	"context"
	"fmt"
	"golek_bookmark_service/pkg/contracts"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	DbUsername            string
	DBPassword            string
	DbName                string
	DbHost                string
	DbPort                string
	DbCollectionBookmarks string
	collection            *mongo.Collection
	connection            *mongo.Database
	config                contracts.DBConfig
}

func New(config contracts.DBConfig) *Database {
	return &Database{
		DbUsername:            config.GetDBConfig()["USERNAME"],
		DBPassword:            config.GetDBConfig()["PASSWORD"],
		DbName:                config.GetDBConfig()["NAME"],
		DbHost:                config.GetDBConfig()["HOST"],
		DbPort:                config.GetDBConfig()["PORT"],
		DbCollectionBookmarks: config.GetDBConfig()["COLLECTION_BOOKMARKS"],
		config:                config,
	}
}

func (db *Database) Prepare() contracts.MongoDBContract {

	//Singleton Connection
	if db.connection == nil {

		clientOptions := options.Client().ApplyURI(db.DSN())

		client, err := mongo.NewClient(clientOptions)
		if err != nil {
			panic(err.Error())
		}

		err = client.Connect(context.Background())
		if err != nil {
			panic(err.Error())
		}

		log.Println("Pinging MongoDB")
		err = client.Ping(context.Background(), nil)
		if err != nil {
			log.Fatalf("Pinging Error %v", err.Error())
		}

		db.connection = client.Database(db.DbName)
		log.Println("Connected to the database MongoDB")

	} else {
		log.Println("Already Connected to the database: MongoDB")
	}

	return db
}

func (db *Database) GetCollection(collection string) *mongo.Collection {

	switch collection {
	case db.DbCollectionBookmarks:
		return db.connection.Collection(collection)
	default:
		return nil
	}
}

func (db *Database) DSN() string {
	// return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?", db.DbUsername, db.DBPassword, db.DbHost, db.DbPort, db.DbName)
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", db.DbUsername, db.DBPassword, db.DbHost, db.DbPort, db.DbName)

}

func (db *Database) GetConnection() *mongo.Database {
	return db.connection
}
