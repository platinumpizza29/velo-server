package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"velo-server/internal/config"
	"velo-server/internal/database"
	"velo-server/internal/handler"
	"velo-server/internal/repository"
	"velo-server/internal/service"
	"velo-server/internal/storage"
)

type Server struct {
	port  int
	db    database.Service
	music *handler.Handler
}

func NewServer() (*http.Server, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db := database.New()
	r2 := storage.NewR2Storage(cfg)
	tracks := repository.NewTrackRepository(db.DB())
	music := handler.NewHandler(service.NewMusicService(tracks, r2))
	app := &Server{
		port:  port,
		db:    db,
		music: music,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}
