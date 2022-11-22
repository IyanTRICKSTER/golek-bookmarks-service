package controllers

import (
	"github.com/gin-gonic/gin"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/http/middleware"
	"golek_bookmark_service/pkg/http/responses"
	"net/http"
)

func SetupHandler(router *gin.Engine, bookmarkUsecase *contracts.BookmarkUsecase) {
	bookmarkHandler := BookmarkHandler{BookmarkUsecase: *bookmarkUsecase}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, responses.HttpResponse{
			StatusCode: http.StatusNotFound,
			Message:    "PAGE NOT FOUND",
			Data:       nil,
		})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, responses.HttpResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "METHOD NOT ALLOWED",
			Data:       nil,
		})
	})

	bRoute := router.Group("/api/bookmark/")
	bRoute.Use(middleware.ValidateRequestHeaderMiddleware)
	bRoute.GET("/", bookmarkHandler.Fetch)
	bRoute.GET("/:id", bookmarkHandler.FetchById)
	bRoute.GET("/u/:user_id", bookmarkHandler.FetchByUserID)
	//bRoute.POST("/create", bookmarkHandler.Create)
	bRoute.DELETE("/course/:user_id", bookmarkHandler.RevokePost)
	bRoute.PATCH("/course/:user_id", bookmarkHandler.AddPost)

}
