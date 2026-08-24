package router

import (
	"movies-api/internal/handlers"
	"net/http"
)

func NewRouter(actorHandler *handlers.ActorHandler,
	genreHandler *handlers.GenreHandler,
	movieHandler *handlers.MovieHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/actor/{id}", actorHandler.GetOneActorHandler)
	mux.HandleFunc("GET /api/actor", actorHandler.GetAllActorsHandler)
	mux.HandleFunc("POST /api/actor", actorHandler.CreateActorHandler)
	mux.HandleFunc("GET /api/actor/{name}", actorHandler.GetActorByNameHandler)
	mux.HandleFunc("DELETE /api/actor/{id}", actorHandler.DeleteActor)
	mux.HandleFunc("PATCH /api/actor/{id}", actorHandler.UpdateActor)

	mux.HandleFunc("GET /api/genre/{id}", genreHandler.GetOneGenreHandler)

	mux.HandleFunc("GET /api/movie/{id}", movieHandler.GetOneMovieHandler)
	mux.HandleFunc("GET /api/movie", movieHandler.GetAllMoviesHAndler)
	mux.HandleFunc("POST /api/movie", movieHandler.CreateMovieHandler)
	mux.HandleFunc("DELETE /api/movie/{id}", movieHandler.DeleteMovieHandler)
	mux.HandleFunc("PATCH /api/movie/{id}", movieHandler.UpdateMovieHandler)
	mux.HandleFunc("GET /api/movie/{name}", movieHandler.GetMovieByNameHandler)

	return mux
}
