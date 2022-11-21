package controllers

import (
	"github.com/gin-gonic/gin"
	"golek_bookmark_service/pkg/contracts"
)

func SetupHandler(router *gin.Engine, bookmarkUsecase *contracts.BookmarkUsecase) {
	bookmarkHandler := BookmarkHandler{BookmarkUsecase: *bookmarkUsecase}

	bRoute := router.Group("/api/bookmark/")
	bRoute.GET("/", bookmarkHandler.Fetch)
	bRoute.GET("/:id", bookmarkHandler.FetchById)
	bRoute.GET("/u/:user_id", bookmarkHandler.FetchByUserID)
	//bRoute.POST("/create", bookmarkHandler.Create)
	bRoute.DELETE("/course/:user_id", bookmarkHandler.RevokePost)
	bRoute.PATCH("/course/:user_id", bookmarkHandler.AddPost)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

}
