package main

import (
	"github.com/gin-gonic/gin"
	"golek_bookmark_service/cmd/grpc_client"
	"golek_bookmark_service/pkg/config"
	"golek_bookmark_service/pkg/database"
	"golek_bookmark_service/pkg/database/migrations"
	"golek_bookmark_service/pkg/http/controllers"
	"golek_bookmark_service/pkg/repositories"
	"golek_bookmark_service/pkg/usecase"
)

func main() {

	engine := gin.Default()

	//Create Config Instance
	cfg := config.New(".env")

	//Connecting Databases
	db := database.New(cfg)
	db.Prepare()

	//Migrations
	mg := migrations.New(db)
	mg.MigrateSettings()

	//Setup Bookmarks
	//Repo
	bookmarkRepo := repositories.NewBookmarkDBRepository(
		db.GetConnection(),
		db.GetCollection(cfg.GetDBConfig()["COLLECTION_BOOKMARKS"]),
	)

	//Connect to Course Service via GRPC
	grpcPostService := grpc_client.New(cfg)
	_, err := grpcPostService.Dial()
	if err != nil {
		panic(err.Error())
	}

	bookmarkUsecase := usecase.NewBookmarkUsecase(bookmarkRepo, grpcPostService)
	//Setup Delivery/Controller
	controllers.SetupHandler(engine, &bookmarkUsecase)

	if port := cfg.GetAppConfig()["PORT"]; port == "" {
		err := engine.Run(":8080")
		if err != nil {
			panic(err)
		}
	} else {
		err := engine.Run(":" + port)
		if err != nil {
			panic(err)
		}
	}

}
