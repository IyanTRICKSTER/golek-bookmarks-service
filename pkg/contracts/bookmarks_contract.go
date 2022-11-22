package contracts

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golek_bookmark_service/pkg/contracts/status"
	"golek_bookmark_service/pkg/http/requests"
	"golek_bookmark_service/pkg/models"
)

type BookmarksRepository interface {
	// Fetch Fetch all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, opStatus status.OperationStatus, err error)
	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error)
	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error)
	Create(ctx context.Context, bookmark *models.Bookmark) (bookmarkID primitive.ObjectID, opStatus status.OperationStatus, err error)
	Update(ctx context.Context, bookmark *models.Bookmark, bookmarkID string) (opStatus status.OperationStatus, err error)
	AddPost(ctx context.Context, userID string, coursesID []string) (opStatus status.OperationStatus, err error)
	Delete(ctx context.Context, bookmarkID string) (opStatus status.OperationStatus, err error)
	RevokePost(ctx context.Context, userID string, coursesID []string) (opStatus status.OperationStatus, err error)
	GenerateModelID() primitive.ObjectID
	GenerateObjectIDFromString(id string) primitive.ObjectID
}

type BookmarkUsecase interface {
	// Fetch Fetch all data from database;
	// 'exclude' param specify which model fields you want to skip/unselect;
	// 'limit' and 'skip param are used to perform some kind of pagination
	Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, opStatus status.OperationStatus, err error)

	// FetchById fetch data by id;
	// 'exclude' param specify which model fields you want to skip/unselect;
	FetchById(ctx context.Context, id string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error)

	FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error)
	Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, opStatus status.OperationStatus, err error)
	AddPost(ctx context.Context, request *requests.AddPostBookmarkRequest, userID string) (opStatus status.OperationStatus, err error)
	RevokePost(ctx context.Context, request *requests.DeleteAttachedPostRequest, userID string) (opStatus status.OperationStatus, err error)
	Delete(ctx context.Context, bookmarkID string) (opStatus status.OperationStatus, err error)
}
