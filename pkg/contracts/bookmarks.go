package contracts

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golek_bookmark_service/pkg/http/requests"
	"golek_bookmark_service/pkg/models"
)

type BookmarksDBRepository interface {
	// Fetch Fetch all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error)
	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error)
	Create(ctx context.Context, bookmark *models.Bookmark) (bookmarkID primitive.ObjectID, err error)
	Update(ctx context.Context, bookmark *models.Bookmark, bookmarkID string) (status bool, err error)
	AddPost(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	Delete(ctx context.Context, bookmarkID string) (status bool, err error)
	RevokePost(ctx context.Context, userID string, coursesID []string) (status bool, err error)
	GenerateModelID() primitive.ObjectID
	GenerateObjectIDFromString(id string) primitive.ObjectID
}

type BookmarkUsecase interface {
	// Fetch Fetch all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error)

	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, err error)

	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error)
	Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, err error)
	AddPost(ctx context.Context, request *requests.AddPostBookmarkRequest, userID string) (status bool, err error)
	RevokePost(ctx context.Context, request *requests.DeleteAttachedPostRequest, userID string) (status bool, err error)
	Delete(ctx context.Context, bookmarkID string) (status bool, err error)
}
