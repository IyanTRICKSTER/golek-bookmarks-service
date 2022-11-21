package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/models"
	"log"
)

type BookmarkDatabaseRepository struct {
	Connection *mongo.Database
	Collection *mongo.Collection
}

func (d BookmarkDatabaseRepository) Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error) {
	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	//Set options
	opts := options.Find()
	opts.SetProjection(excluded)
	opts.SetLimit(limit)
	opts.SetSkip(skip)

	//Fetch Records
	filter := map[string]interface{}{"deleted_at": nil}

	records, err := d.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	//Close Cursor
	defer func(records *mongo.Cursor, ctx context.Context) {
		err := records.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}(records, ctx)

	bookmarks = make([]models.Bookmark, 0)

	err = records.All(ctx, &bookmarks)
	if err != nil {
		return nil, err
	}

	return bookmarks, nil

}

func (d BookmarkDatabaseRepository) FetchById(ctx context.Context, id string, exclude []string) (bookmarks models.Bookmark, err error) {

	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	//Set options
	opts := options.FindOne().SetProjection(excluded)

	var bookmark models.Bookmark
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return bookmark, err
	}

	filter := map[string]interface{}{"_id": objectID, "deleted_at": nil}
	err = d.Collection.FindOne(ctx, filter, opts).Decode(&bookmark)

	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

func (d BookmarkDatabaseRepository) FetchByUserId(ctx context.Context, userId string, exclude []string) (bookmarks models.Bookmark, err error) {

	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	opts := options.FindOne().SetProjection(excluded)

	var bookmark models.Bookmark

	filter := map[string]interface{}{"user_id": userId, "deleted_at": nil}

	err = d.Collection.FindOne(ctx, filter, opts).Decode(&bookmark)
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil
}

func (d BookmarkDatabaseRepository) Create(ctx context.Context, bookmark *models.Bookmark) (postID primitive.ObjectID, err error) {

	//	Use Transaction
	err = d.Connection.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {

		// Start Transaction
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		// Insert Data To the Database & abort if it fails
		insertedData, err := d.Collection.InsertOne(ctx, bookmark)
		if err != nil {
			err := sessionContext.AbortTransaction(ctx)
			if err != nil {
				return err
			}
			return err
		}

		postID = insertedData.InsertedID.(primitive.ObjectID)

		// Commit Data if no error
		err = sessionContext.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return postID, nil
}

func (d BookmarkDatabaseRepository) Update(ctx context.Context, bookmark *models.Bookmark, bookmarkID string) (status bool, err error) {

	objectId, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return false, err
	}

	filter := bson.D{{"_id", objectId}}
	_, err = d.Collection.UpdateOne(ctx, filter, bson.D{{"$set", bookmark}})
	if err != nil {
		return false, err
	}
	return true, err
}

func (d BookmarkDatabaseRepository) Delete(ctx context.Context, bookmarkID string) (status bool, err error) {

	objectID, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return false, err
	}

	_, err = d.Collection.DeleteOne(ctx, bson.D{{"_id", objectID}})
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d BookmarkDatabaseRepository) AddPost(ctx context.Context, userID string, postIDs []string) (status bool, err error) {

	filter := bson.D{{"user_id", userID}}

	postsObjID := make([]bson.D, 0)
	for _, c := range postIDs {
		postsObjID = append(postsObjID, bson.D{{"id", d.GenerateObjectIDFromString(c)}})
	}

	statement := bson.M{"$addToSet": bson.M{"posts": bson.M{"$each": postsObjID}}}

	result, err := d.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY ADD POST: ", err.Error())
		return false, err
	}

	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY ADD POST: document not matched")
		return false, errors.New("document not matched")

	}

	if result.ModifiedCount == 0 {
		log.Println("BOOKMARK REPOSITORY ADD POST: document not modified")
		return false, errors.New("document not modified")
	}

	return true, nil
}

func (d BookmarkDatabaseRepository) RevokePost(ctx context.Context, userID string, postIDs []string) (status bool, err error) {

	filter := bson.D{{"user_id", userID}}

	cID := make([]primitive.ObjectID, 0)
	for _, s := range postIDs {
		cID = append(cID, d.GenerateObjectIDFromString(s))
	}

	statement := bson.M{"$pull": bson.M{"posts": bson.M{"id": bson.M{"$in": cID}}}}

	result, err := d.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY DELETE POST: ", err.Error())
		return false, err
	}

	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY DELETE POST: document not matched")
		return false, errors.New("document not matched")

	}

	if result.ModifiedCount == 0 {
		log.Println("BOOKMARK REPOSITORY DELETE POST: document not modified")
		//return false, errors.New("document not modified")
	}

	return true, nil
}

func (d BookmarkDatabaseRepository) GenerateModelID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (d BookmarkDatabaseRepository) GenerateObjectIDFromString(id string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return hex
}

func NewBookmarkDBRepository(conn *mongo.Database, coll *mongo.Collection) contracts.BookmarksDBRepository {

	return &BookmarkDatabaseRepository{
		Connection: conn,
		Collection: coll,
	}
}
