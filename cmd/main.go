package main

import (
	"log"
	"net/http"

	"movies-api/internal/database"
	"movies-api/internal/handlers"
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

	mux := router.NewRouter(actorHandler,
		genreHandler,
		movieHandler,
	)

	log.Println("\033[32mListening on http://localhost:8080\033[0m")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
