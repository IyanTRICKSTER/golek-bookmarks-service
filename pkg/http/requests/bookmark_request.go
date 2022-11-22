package requests

type CreateBookmarkRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Posts  []Post `json:"posts" binding:"required,dive"`
}

type AddPostBookmarkRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Posts  []Post `json:"posts" binding:"required,dive"`
}

type DeleteAttachedPostRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Posts  []Post `json:"posts" binding:"required,dive"`
}

type Post struct {
	ID string `json:"id" binding:"required"`
}
