package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/contracts/status"
	"golek_bookmark_service/pkg/http/requests"
	"golek_bookmark_service/pkg/http/responses"
	"golek_bookmark_service/pkg/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type BookmarkHandler struct {
	BookmarkUsecase contracts.BookmarkUsecase
}

func (h *BookmarkHandler) Fetch(c *gin.Context) {

	var excludedField []string
	if c.Query("exclude") != "" {
		excludedField = strings.Split(c.Query("exclude"), ",")
	}

	page, ok := c.GetQuery("page")
	if page == "" || !ok {
		page = "1"
	} else if page == "0" {
		page = "1"
	}

	qPage, err2 := strconv.ParseInt(page, 10, 64)
	if err2 != nil {
		return
	}

	paginate := models.Pagination{
		Page:    qPage,
		PerPage: 25,
	}

	limit, skip := paginate.GetPagination()
	bookmarks, _, err := h.BookmarkUsecase.Fetch(c.Request.Context(), excludedField, limit, skip)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.HttpPaginationResponse{
		PerPage: paginate.PerPage,
		Page:    paginate.Page,
		HttpResponse: responses.HttpResponse{
			Data:       bookmarks,
			StatusCode: http.StatusOK,
		},
	})
	return
}

func (h BookmarkHandler) FetchById(c *gin.Context) {

	bookmark, opStatus, err := h.BookmarkUsecase.FetchById(c.Request.Context(), c.Param("id"), []string{})
	if err != nil {
		//return 404 not found
		if status.Is(opStatus, status.BookmarkNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
	return
}

func (h BookmarkHandler) FetchByUserID(c *gin.Context) {
	bookmark, opStatus, err := h.BookmarkUsecase.FetchByUserId(c.Request.Context(), c.Param("user_id"), []string{})
	if err != nil {
		//return 404 not found
		if status.Is(opStatus, status.BookmarkNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
}

func (h BookmarkHandler) Create(c *gin.Context) {

	val, _ := c.Get("authenticatedRequest")
	authContext := context.WithValue(context.Background(), "authenticatedRequest", val)

	var createRequest requests.CreateBookmarkRequest

	err := c.ShouldBindJSON(&createRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	bookmark, opStatus, err := h.BookmarkUsecase.Create(authContext, &createRequest)
	if err != nil {
		if status.Is(opStatus, status.BookmarkDuplicationOccurs) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if status.Is(opStatus, status.OperationUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookmark)
	return
}

func (h BookmarkHandler) AddPost(c *gin.Context) {

	val, _ := c.Get("authenticatedRequest")
	authContext := context.WithValue(context.Background(), "authenticatedRequest", val)

	var addPostReq requests.AddPostBookmarkRequest

	err := c.ShouldBindJSON(&addPostReq)
	if err != nil {
		log.Println("BOOKMARK HANDLER: AddPost", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opStatus, err := h.BookmarkUsecase.AddPost(authContext, &addPostReq, c.Param("user_id"))
	if err != nil || status.Is(opStatus, status.BookmarkPostFailed) {
		//return 404 not found
		if status.Is(opStatus, status.BookmarkNotExist) {
			log.Println("BOOKMARK HANDLER: AddPost", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if status.Is(opStatus, status.OperationUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if status.Is(opStatus, status.OperationForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		log.Println("BOOKMARK HANDLER: AddPost", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return

}

func (h BookmarkHandler) RevokePost(c *gin.Context) {

	val, _ := c.Get("authenticatedRequest")
	authContext := context.WithValue(context.Background(), "authenticatedRequest", val)

	var revokePostReq requests.DeleteAttachedPostRequest

	err := c.ShouldBindJSON(&revokePostReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opStatus, err := h.BookmarkUsecase.RevokePost(authContext, &revokePostReq, c.Param("user_id"))
	if err != nil || status.Is(opStatus, status.BookmarkPostFailed) {
		//return 404 not found
		if status.Is(opStatus, status.BookmarkNotExist) {
			log.Println("BOOKMARK HANDLER: RevokePost", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if status.Is(opStatus, status.OperationUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if status.Is(opStatus, status.OperationForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	return
}
