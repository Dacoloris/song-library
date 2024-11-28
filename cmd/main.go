package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"song-library/configs"
	"song-library/internal/handler"
	"song-library/internal/models"
	"song-library/internal/repository"
	"song-library/internal/service"
	"song-library/pkg/db"
	"song-library/pkg/logger"

	_ "song-library/cmd/docs"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := configs.LoadConfig()

	db := db.NewDb(cfg.Dsn)
	err = db.DB.AutoMigrate(&models.Song{})
	if err != nil {
		panic(err)
	}

	songRepo := repository.NewSongRepository(db)
	songService := service.NewSongService(songRepo)
	router := gin.Default()
	handler.NewSongHandler(router, songService, log, cfg.ApiUrl)

	log.Info("Server is running on ", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Error("failed to start server", zap.Error(err))
	}
}
