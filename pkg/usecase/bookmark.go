package usecase

import (
	"context"
	"errors"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/http/requests"
	"golek_bookmark_service/pkg/models"
	"log"
	"time"
)

type BookmarkUsecase struct {
	DBRepository          contracts.BookmarksDBRepository
	GRPCPostServiceClient contracts.GRPCPostService
}

func (b BookmarkUsecase) Fetch(ctx context.Context, exclude []string, limit int64, skip int64) (bookmarks []models.Bookmark, err error) {
	bookmarks, err = b.DBRepository.Fetch(ctx, exclude, limit, skip)
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (b BookmarkUsecase) FetchById(ctx context.Context, bookmarkID string, exclude []string) (bookmark models.Bookmark, err error) {

	//Fetch Bookmark Containing embedded post id
	bookmark, err = b.DBRepository.FetchById(ctx, bookmarkID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchById ERROR", err)
		return models.Bookmark{}, err
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
	bookmark.Posts = posts

	//log.Println(bookmark)
	return bookmark, nil
}

func (b BookmarkUsecase) FetchByUserId(ctx context.Context, userID string, exclude []string) (bookmark models.Bookmark, err error) {
	bookmark, err = b.DBRepository.FetchByUserId(ctx, userID, exclude)
	if err != nil {
		log.Println("BOOKMARK USECASE: FetchByUserId ERROR >>", err)
		return models.Bookmark{}, err
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
	bookmark.Posts = posts

	return bookmark, nil
}

func (b BookmarkUsecase) Create(ctx context.Context, request *requests.CreateBookmarkRequest) (bookmark models.Bookmark, err error) {

	posts := make([]models.Post, 0)
	for _, post := range request.Posts {
		posts = append(posts, models.Post{ID: b.DBRepository.GenerateObjectIDFromString(post.ID)})
	}

	timeNow := time.Now()
	newBookmark := models.Bookmark{
		ID:        b.DBRepository.GenerateModelID(),
		UserID:    request.UserID,
		Posts:     posts,
		UpdatedAt: &timeNow,
		CreatedAt: &timeNow,
	}

	bookmarkID, err := b.DBRepository.Create(ctx, &newBookmark)
	if err != nil {
		return models.Bookmark{}, err
	}

	newBookmark.ID = bookmarkID
	return newBookmark, nil
}

func (b BookmarkUsecase) AddPost(ctx context.Context, request *requests.AddPostBookmarkRequest, userID string) (status bool, err error) {

	//if a bookmark not found, then create a new one
	_, err = b.FetchByUserId(ctx, userID, []string{})
	if err != nil {
		if err.Error() == errors.New("mongo: no documents in result").Error() {
			_, err = b.Create(ctx, (*requests.CreateBookmarkRequest)(request))
			if err != nil {
				log.Println("BOOKMARK USECASE: AddPost >>", err)
				return false, err
			}
			log.Println("BOOKMARK USECASE: AddPost >>", "Post has been added", request.Posts)
			return true, nil
		}
		log.Println("BOOKMARK USECASE: AddPost >>", err)
		return false, err
	}

	postID := make([]string, 0)
	for _, post := range request.Posts {
		postID = append(postID, post.ID)
	}

	status, err = b.DBRepository.AddPost(ctx, userID, postID)
	if err != nil {
		log.Println("BOOKMARK USECASE: AddPost >>", err)
		return false, err
	}

	return status, nil
}

func (b BookmarkUsecase) RevokePost(ctx context.Context, request *requests.DeleteAttachedPostRequest, userID string) (status bool, err error) {

	postsID := make([]string, 0)
	for _, c := range request.Posts {
		postsID = append(postsID, c.ID)
	}

	status, err = b.DBRepository.RevokePost(ctx, userID, postsID)
	if err != nil {
		log.Println("BOOKMARK USECASE REVOKE post:", err.Error())
		return false, err
	}

	return status, nil
}

func (b BookmarkUsecase) Delete(ctx context.Context, bookmarkID string) (status bool, err error) {
	status, err = b.DBRepository.Delete(ctx, bookmarkID)
	if err != nil {
		return false, err
	}
	return status, nil
}

func NewBookmarkUsecase(DBRepository contracts.BookmarksDBRepository, GRPCPostServiceClient contracts.GRPCPostService) contracts.BookmarkUsecase {
	return &BookmarkUsecase{DBRepository: DBRepository, GRPCPostServiceClient: GRPCPostServiceClient}
}
