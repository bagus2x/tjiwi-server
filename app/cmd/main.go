package main

import (
	"log"

	"github.com/bagus2x/tjiwi/app/handler"
	appMiddleware "github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	basepaperrepo "github.com/bagus2x/tjiwi/pkg/basepaper/repository"
	basepaperservice "github.com/bagus2x/tjiwi/pkg/basepaper/service"
	historyrepo "github.com/bagus2x/tjiwi/pkg/history/repository"
	historyservice "github.com/bagus2x/tjiwi/pkg/history/service"
	storageRepo "github.com/bagus2x/tjiwi/pkg/storage/repository"
	storageService "github.com/bagus2x/tjiwi/pkg/storage/service"
	stormembRepo "github.com/bagus2x/tjiwi/pkg/storagemember/repository"
	stormembService "github.com/bagus2x/tjiwi/pkg/storagemember/service"
	userrepo "github.com/bagus2x/tjiwi/pkg/user/repository"
	userservice "github.com/bagus2x/tjiwi/pkg/user/service"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.New()
	cfg := config.New()

	database := db.NewPostgresDatabase(cfg)

	userRepo := userrepo.New(database)
	storageRepo := storageRepo.New(database)
	stormembRepo := stormembRepo.New(database)
	basePaperRepo := basepaperrepo.New(database)
	historyRepo := historyrepo.New(database)

	userService := userservice.New(userRepo, cfg)
	stormembService := stormembService.New(stormembRepo, cfg)
	storageService := storageService.New(storageRepo, stormembRepo, cfg)
	basePaperService := basepaperservice.New(basePaperRepo, historyRepo)
	historyService := historyservice.New(historyRepo)

	mw := appMiddleware.New(userService, stormembService)

	app.Use(gin.Recovery())
	app.Use(gin.Logger())
	app.Use(mw.Cors())
	app.Use(mw.GinContextToContextMiddleware())

	userGroup := app.Group("/users")
	storageGroup := app.Group("/storages")
	stormembGroup := app.Group("/storagemembers")
	basePaper := app.Group("/basepapers")
	history := app.Group("/histories")

	handler.User(userGroup, userService, mw)
	handler.Storage(storageGroup, storageService, mw)
	handler.StorageMember(stormembGroup, stormembService, mw)
	handler.BasePaper(basePaper, basePaperService, mw)
	handler.History(history, historyService, mw)

	log.Fatal(app.Run(cfg.AppPort()))
}
