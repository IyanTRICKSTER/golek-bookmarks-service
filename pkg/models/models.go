package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type Pagination struct {
	Page    int64
	PerPage int64
}

func (p Pagination) GetPagination() (limit int64, skip int64) {
	return p.PerPage, (p.Page - 1) * p.PerPage
}

func GenerateObjectIDFromHex(hex string) primitive.ObjectID {
	objectID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println("MODELS GENERATE OBJECT ID:", err.Error())
		return primitive.NilObjectID
	}

	return objectID
}

func GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}
