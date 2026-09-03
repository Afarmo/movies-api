package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movies-api/internal/database"
	"movies-api/internal/handlers"
	"movies-api/internal/middleware"
	"movies-api/internal/repository"
	"movies-api/internal/router"
	"movies-api/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewDataBase("movieapi.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.InitializeQuery(db); err != nil {
		log.Fatal(err)
	}

	actorRepo := repository.NewActorRepo(db)
	genreRepo := repository.NewGenreRepo(db)
	movieRepo := repository.NewMovieRepo(db)

	actorService := service.NewActorService(actorRepo)
	genreService := service.NewGenreService(genreRepo)
	movieService := service.NewMovieService(movieRepo)

	actorHandler := handlers.NewActorHandler(actorService)
	genreHandler := handlers.NewGenreHandler(genreService)
	movieHandler := handlers.NewMovieHandler(movieService)

	mux := router.NewRouter(actorHandler, genreHandler, movieHandler)
	handler := middleware.Recover(middleware.Logger(mux))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println("\033[36m  [STARTUP]\033[0m  Listening on http://localhost:8080")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("\033[31m[PANIC]\033[0m Server error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("\033[34m[SHUTDOWN]\033[0m Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("\033[31m[PANIC]\033[0m Forced shutdown:", err)
	}

	log.Println("\033[34m  [SHUTDOWN]\033[0m Server stopped")
}
