package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Bookmark struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Posts     []Post             `json:"posts" bson:"posts"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedAt *time.Time         `json:"created_at,omitempty" bson:"created_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at"`
}

type Post struct {
	ID   primitive.ObjectID `json:"id" bson:"id"`
	Name string             `json:"name,omitempty" bson:"-"`
}
