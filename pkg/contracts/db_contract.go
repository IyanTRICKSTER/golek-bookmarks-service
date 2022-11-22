package contracts

import "go.mongodb.org/mongo-driver/mongo"

type DBContract interface {
	DSN() string
	GetConnection() *mongo.Database
}

type MongoDBContract interface {
	GetCollection(collection string) *mongo.Collection
	DBContract
}
