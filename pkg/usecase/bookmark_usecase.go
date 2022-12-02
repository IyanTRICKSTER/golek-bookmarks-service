package usecase

import (
	"context"
	"errors"
	"fmt"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/contracts/status"
	"golek_bookmark_service/pkg/http/middleware"
	"golek_bookmark_service/pkg/http/requests"
	"golek_bookmark_service/pkg/models"
	"log"
	"strings"
	"time"
)

type BookmarkUsecase struct {
	DBRepository          contracts.BookmarksRepository
	GRPCPostServiceClient contracts.GRPCPostService
}

func NewBookmarkUsecase(DBRepository contracts.BookmarksRepository, GRPCPostServiceClient contracts.GRPCPostService) contracts.BookmarkUsecase {
	return &BookmarkUsecase{DBRepository: DBRepository, GRPCPostServiceClient: GRPCPostServiceClient}
}

func (b BookmarkUsecase) Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, opStatus status.OperationStatus, err error) {

	bookmarks, opStatus, err = b.DBRepository.Fetch(ctx, exclude, limit, skip)
	if err != nil {
		return nil, opStatus, err
	}
	return bookmarks, opStatus, nil
}

func (b BookmarkUsecase) FetchById(ctx context.Context, bookmarkID string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error) {

	//Fetch Bookmark Containing embedded post id
	bookmark, opStatus, err = b.DBRepository.FetchById(ctx, bookmarkID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchById ERROR", err)
		return bookmark, opStatus, err
	}

	//Fetch Post Data From postService through GRPC
	pIDs := make([]string, 0)
	for _, c := range bookmark.Posts {
		pIDs = append(pIDs, c.ID.Hex())
	}

	posts, err := b.GRPCPostServiceClient.Fetch(ctx, pIDs)
	if err != nil {
		log.Printf("Fetching Bookmark Failed >> %v", err.Error())
	}
	//log.Println("BOOKMARK USECASE: FETCH BY ID: gRPC postService Result >>", posts)

	//Attach post data from postService to a Bookmark
	if err == nil {
		if len(posts) != 0 {
			bookmark.Posts = posts
		}
	}

	//log.Println(bookmark)
	return bookmark, opStatus, nil
}

func (b BookmarkUsecase) FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, opStatus status.OperationStatus, err error) {

	bookmark, opStatus, err = b.DBRepository.FetchByUserId(ctx, userID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchByUserId ERROR >>", err)
		return bookmark, opStatus, err
	}

	//Fetch Post Data From postService through GRPC
	pIDs := make([]string, 0)
	for _, c := range bookmark.Posts {
		pIDs = append(pIDs, c.ID.Hex())
	}

	//Attach post data from postService to a Bookmark
	posts, err := b.GRPCPostServiceClient.Fetch(ctx, pIDs)
	if err != nil {
		log.Printf("Fetching Bookmark Failed >> %v", err.Error())
	}
	//Attach post data from postService to a Bookmark
	if err == nil {
		if len(posts) != 0 {
			bookmark.Posts = posts
		}
	}

	return bookmark, opStatus, nil
}

func (b BookmarkUsecase) Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, opStatus status.OperationStatus, err error) {

	//Check user authorization
	authenticated, opStatus, err := ProtectResource(ctx, contracts.Resource{
		Alias: "c",
		Name:  "Create Service",
	}, models.Bookmark{}, func(isOwner bool) (opStatus status.OperationStatus, err error) {
		return status.OperationAuthorized, nil
	})
	if err != nil {
		return models.Bookmark{}, opStatus, err
	}

	if authenticated.UserID != request.UserID {
		return models.Bookmark{}, status.OperationForbidden, errors.New("user id doesn't match with authenticated token")
	}

	posts := make([]models.Post, 0)
	for _, post := range request.Posts {
		posts = append(posts, models.Post{ID: b.DBRepository.GenerateObjectIDFromString(post.ID)})
	}

	timeNow := time.Now()
	newBookmark := models.Bookmark{
		ID:        b.DBRepository.GenerateModelID(),
		UserID:    authenticated.UserID,
		Posts:     posts,
		UpdatedAt: &timeNow,
		CreatedAt: &timeNow,
	}

	bookmarkID, opStatus, err := b.DBRepository.Create(ctx, &newBookmark)
	if err != nil {
		return models.Bookmark{}, opStatus, err
	}

	newBookmark.ID = bookmarkID
	return newBookmark, opStatus, nil
}

func (b BookmarkUsecase) AddPost(ctx context.Context, request *requests.AddPostBookmarkRequest, userID string) (opStatus status.OperationStatus, err error) {

	//if a bookmark not found, then create a new one
	bookmark, opStatus, err := b.DBRepository.FetchByUserId(ctx, userID, []string{})
	if err != nil {

		//if bookmark not exists, create new
		if opStatus == status.BookmarkNotExist {
			_, opStatus, err = b.Create(ctx, (*requests.CreateBookmarkRequest)(request))
			if err != nil {
				if opStatus == status.OperationUnauthorized {
					log.Println("BOOKMARK USECASE: AddPost >>", err.Error())
					return status.OperationUnauthorized, err
				}
				log.Println("BOOKMARK USECASE: AddPost >>", err.Error())
				return status.BookmarkPostFailed, err
			}
			log.Println("BOOKMARK USECASE: AddPost >>", "Post has been added", request.Posts)
			return status.BookmarkPostSuccess, nil
		}

		log.Println("BOOKMARK USECASE: AddPost >>", err)
		return status.BookmarkPostFailed, err
	}

	//Check user authorization & model owner
	authenticated, opStatus, err := ProtectResource(ctx, contracts.Resource{
		Alias: "u",
		Name:  "Add Post Service",
	}, bookmark, func(isOwner bool) (opStatus status.OperationStatus, err error) {
		if !isOwner {
			return status.OperationForbidden, errors.New("you are not the owner")
		}
		return status.OperationAuthorized, nil
	})
	if err != nil {
		return opStatus, err
	}

	if authenticated.UserID != request.UserID && authenticated.UserID != userID {
		return status.OperationForbidden, errors.New("user id doesn't match with authenticated token")
	}

	postID := make([]string, 0)
	for _, post := range request.Posts {
		postID = append(postID, post.ID)
	}

	opStatus, err = b.DBRepository.AddPost(ctx, userID, postID)
	if err != nil {
		log.Println("BOOKMARK USECASE: AddPost >>", err)
		return opStatus, err
	}

	return opStatus, nil
}

func (b BookmarkUsecase) RevokePost(ctx context.Context, request *requests.DeleteAttachedPostRequest, userID string) (opStatus status.OperationStatus, err error) {

	bookmark, opStatus, err := b.DBRepository.FetchByUserId(ctx, userID, []string{})
	if err != nil {
		//if bookmark not exists, create new
		if opStatus == status.BookmarkNotExist {
			return status.BookmarkNotExist, err
		}
		return status.BookmarkPostRevokeFailed, err
	}

	//Check user authorization & model owner
	authenticated, opStatus, err := ProtectResource(ctx, contracts.Resource{
		Alias: "u",
		Name:  "Revoke Post Service",
	}, bookmark, func(isOwner bool) (opStatus status.OperationStatus, err error) {
		if !isOwner {
			return status.OperationForbidden, errors.New("you are not the owner")
		}
		return status.OperationAuthorized, nil
	})
	if err != nil {
		return opStatus, err
	}

	if authenticated.UserID != request.UserID && authenticated.UserID != userID {
		return status.OperationForbidden, errors.New("user id doesn't match with authenticated token")
	}

	postsID := make([]string, 0)
	for _, c := range request.Posts {
		postsID = append(postsID, c.ID)
	}

	opStatus, err = b.DBRepository.RevokePost(ctx, userID, postsID)
	if err != nil {
		log.Println("BOOKMARK USECASE REVOKE post:", err.Error())
		return opStatus, err
	}

	return opStatus, nil
}

func (b BookmarkUsecase) Delete(ctx context.Context, bookmarkID string) (opStatus status.OperationStatus, err error) {

	opStatus, err = b.DBRepository.Delete(ctx, bookmarkID)
	if err != nil {
		return opStatus, err
	}
	return opStatus, nil
}

// ProtectResource Test
func ProtectResource(ctx context.Context, resource contracts.Resource, model models.Bookmark, callback func(isOwner bool) (opStatus status.OperationStatus, err error)) (*middleware.AuthenticatedRequest, status.OperationStatus, error) {

	log.Println("Checking User Permissions")

	//Check User Authorization
	authenticated := ctx.Value("authenticatedRequest").(*middleware.AuthenticatedRequest)
	if !strings.Contains(authenticated.Permissions, resource.Alias) {
		return authenticated, status.OperationUnauthorized, errors.New(fmt.Sprintf("User ID %v doesn't have any permission to access %v resource", authenticated.UserID, resource.Name))
	}

	//Check Model's owner
	isOwner := authenticated.UserID == model.UserID

	//Run the callback
	opStatus, err := callback(isOwner)
	if err != nil {
		return authenticated, opStatus, err
	}

	return authenticated, opStatus, nil
}
