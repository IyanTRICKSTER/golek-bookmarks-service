package migrations

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func (m Migration) MigrateSettings() {
	m.CreateIndexes()
	log.Println("Migrates Settings Success")
}

func (m Migration) CreateIndexes() {

	//_, err := m.DB.GetCollection(m.DB.DbCollectionBookmarks).Indexes().DropOne(context.Background(), "user_id_1")
	//if err != nil {
	//	panic(err)
	//}

	//set bookmark user id as unique value,
	_, err := m.DB.GetCollection(m.DB.DbCollectionBookmarks).Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
	if err != nil {
		log.Println(err)
	}

	_, err = m.DB.GetCollection(m.DB.DbCollectionBookmarks).Indexes().CreateOne(context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "courses.id", Value: 1}},
			Options: options.Index().SetUnique(false),
		})
	if err != nil {
		log.Println(err)

	}
}
