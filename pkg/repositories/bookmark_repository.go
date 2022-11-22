package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/contracts/status"
	"golek_bookmark_service/pkg/models"
	"log"
)

type BookmarkRepository struct {
	Connection *mongo.Database
	Collection *mongo.Collection
}

func (d BookmarkRepository) Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, opStatus status.OperationStatus, err error) {

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
		return nil, status.BookmarkFetchingFailed, err
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
		return nil, status.BookmarkFetchingFailed, err
	}

	return bookmarks, status.OperationSuccess, nil

}

func (d BookmarkRepository) FetchById(ctx context.Context, id string, exclude []string) (bookmarks models.Bookmark, opStatus status.OperationStatus, err error) {

	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	var bookmark models.Bookmark

	//Convert model id from string to mongodb objectID
	modelID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return bookmark, status.BookmarkFetchingFailed, err
	}

	//Set options
	opts := options.FindOne().SetProjection(excluded)
	//Set filters
	filter := map[string]interface{}{"_id": modelID, "deleted_at": nil}
	err = d.Collection.FindOne(ctx, filter, opts).Decode(&bookmark)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return bookmark, status.BookmarkNotExist, err
		}
		return bookmark, status.BookmarkFetchingFailed, err
	}

	return bookmark, status.OperationSuccess, nil
}

func (d BookmarkRepository) FetchByUserId(ctx context.Context, userId string, exclude []string) (bookmarks models.Bookmark, opStatus status.OperationStatus, err error) {

	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	//set options
	opts := options.FindOne().SetProjection(excluded)

	//set filters
	filter := map[string]interface{}{"user_id": userId, "deleted_at": nil}

	var bookmark models.Bookmark
	err = d.Collection.FindOne(ctx, filter, opts).Decode(&bookmark)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return bookmark, status.BookmarkNotExist, err
		}
		return bookmark, status.BookmarkFetchingFailed, err
	}

	return bookmark, status.OperationSuccess, nil
}

func (d BookmarkRepository) Create(ctx context.Context, bookmark *models.Bookmark) (postID primitive.ObjectID, opStatus status.OperationStatus, err error) {

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
		if mongo.IsDuplicateKeyError(err) {
			return primitive.NilObjectID, status.BookmarkDuplicationOccurs, err
		}
		return primitive.NilObjectID, status.BookmarkCreateFailed, err
	}

	return postID, status.BookmarkCreateSuccess, nil
}

func (d BookmarkRepository) Update(ctx context.Context, bookmark *models.Bookmark, bookmarkID string) (opStatus status.OperationStatus, err error) {

	objectId, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return status.BookmarkUpdateFailed, err
	}

	filter := bson.D{{"_id", objectId}}
	result, err := d.Collection.UpdateOne(ctx, filter, bson.D{{"$set", bookmark}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return status.BookmarkNotExist, err
		}
		return status.BookmarkUpdateFailed, err
	}
	if result.MatchedCount != 0 {
		return status.BookmarkUpdateSuccess, nil
	}
	return status.BookmarkUpdateFailed, err
}

func (d BookmarkRepository) Delete(ctx context.Context, bookmarkID string) (opStatus status.OperationStatus, err error) {

	objectID, err := primitive.ObjectIDFromHex(bookmarkID)
	if err != nil {
		return status.BookmarkDeleteFailed, err
	}

	//set filters
	filter := bson.D{{"_id", objectID}}

	result, err := d.Collection.DeleteOne(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return status.BookmarkNotExist, err
		}
		return status.BookmarkDeleteFailed, err
	}

	if result.DeletedCount == 0 {
		return status.BookmarkNotExist, err
	}

	return status.BookmarkDeleteSuccess, nil
}

func (d BookmarkRepository) AddPost(ctx context.Context, userID string, postIDs []string) (opStatus status.OperationStatus, err error) {

	//set filters
	//1. Query by user id
	filter := bson.D{{"user_id", userID}}

	//Convert postIDs string to ObjectID
	postObjIDs := make([]bson.D, 0)
	for _, c := range postIDs {
		postObjIDs = append(postObjIDs, bson.D{{"id", d.GenerateObjectIDFromString(c)}})
	}

	//Set statements
	//1. Append post id to posts array field, if already exists in array, id will not be added
	statement := bson.M{"$addToSet": bson.M{"posts": bson.M{"$each": postObjIDs}}}

	//execute
	result, err := d.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY ADD POST: ", err.Error())
		return status.BookmarkPostFailed, err
	}

	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY ADD POST: document not matched")
		return status.BookmarkNotExist, errors.New("document not matched")
	}

	//if result.ModifiedCount == 0 {
	//	log.Println("BOOKMARK REPOSITORY ADD POST: document not modified")
	//	return status.BookmarkMatchedNotModified, errors.New("document not modified")
	//}

	return status.BookmarkPostSuccess, nil
}

func (d BookmarkRepository) RevokePost(ctx context.Context, userID string, postIDs []string) (opStatus status.OperationStatus, err error) {

	//set filters
	//1. Query by user id
	filter := bson.D{{"user_id", userID}}

	//Convert postIDs string to ObjectID
	cID := make([]primitive.ObjectID, 0)
	for _, s := range postIDs {
		cID = append(cID, d.GenerateObjectIDFromString(s))
	}

	//Set statements
	//1. remove post id matched in postObjIDs
	statement := bson.M{"$pull": bson.M{"posts": bson.M{"id": bson.M{"$in": cID}}}}

	//execute
	result, err := d.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY DELETE POST: ", err.Error())
		return status.BookmarkDeletePostFailed, err
	}

	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY DELETE POST: document not matched")
		return status.BookmarkNotExist, errors.New("document not matched")

	}

	//if result.ModifiedCount == 0 {
	//	log.Println("BOOKMARK REPOSITORY DELETE POST: document not modified")
	//	//return false, errors.New("document not modified")
	//}

	return status.BookmarkDeletePostSuccess, nil
}

func (d BookmarkRepository) GenerateModelID() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (d BookmarkRepository) GenerateObjectIDFromString(id string) primitive.ObjectID {
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return hex
}

func NewBookmarkDBRepository(conn *mongo.Database, coll *mongo.Collection) contracts.BookmarksRepository {

	return &BookmarkRepository{
		Connection: conn,
		Collection: coll,
	}
}
