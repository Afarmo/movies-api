package main

import (
	"log"
	"movies-api/internal/database"
	"movies-api/internal/handlers"
	"movies-api/internal/repository"
	"movies-api/internal/router"
	"movies-api/internal/service"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := database.NewDataBase("movieapi.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := database.InitializeQuery(db); err != nil {
		log.Fatal(err)
	}

	ActorRepo := repository.NewActorRepo(db)
	GenreRepo := repository.NewGenreRepo(db)
	MovieRepo := repository.NewMovieRepo(db)

	ActorService := service.NewActorService(*ActorRepo)
	GenreService := service.NewGenreService(*GenreRepo)
	MovieService := service.NewMovieService(*MovieRepo)

	ActorHandler := handlers.NewActorHandler(ActorService)
	GenreHandler := handlers.NewGenreHandler(GenreService)
	MovieHandler := handlers.NewMovieHandler(MovieService)

	router := router.NewRouter(ActorHandler,
		GenreHandler,
		MovieHandler,
	)

	log.Println("\033[32mListening on http://localhost:8080\033[0m")
	log.Fatal(http.ListenAndServe(":8080", router))
}
