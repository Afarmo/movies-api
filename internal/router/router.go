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

	mux.HandleFunc("GET /api/actors", actorHandler.GetAllActorsHandler)
	mux.HandleFunc("GET /api/actors/search", actorHandler.SearchActorsHandler)
	mux.HandleFunc("GET /api/actors/{id}", actorHandler.GetOneActorHandler)
	mux.HandleFunc("GET /api/movies/{movieId}/actors", actorHandler.GetActorsByMovie)
	mux.HandleFunc("POST /api/actors", actorHandler.CreateActorHandler)
	mux.HandleFunc("DELETE /api/actors/{id}", actorHandler.DeleteActorHandler)
	mux.HandleFunc("PATCH /api/actors/{id}", actorHandler.UpdateActorHandler)

	mux.HandleFunc("GET /api/genres", genreHandler.GetAllGenresHandler)
	mux.HandleFunc("GET /api/genres/{id}", genreHandler.GetOneGenreHandler)
	mux.HandleFunc("POST /api/genres", genreHandler.CreateGenreHandler)
	mux.HandleFunc("DELETE /api/genres/{id}", genreHandler.DeleteGenreHandler)
	mux.HandleFunc("PATCH /api/genres/{id}", genreHandler.UpdateGenreHandler)

	mux.HandleFunc("GET /api/movies", movieHandler.GetAllMoviesHandler)
	mux.HandleFunc("GET /api/movies/search", movieHandler.SearchMoviesHandler)
	mux.HandleFunc("GET /api/movies/{id}", movieHandler.GetOneMovieHandler)
	mux.HandleFunc("POST /api/movies", movieHandler.CreateMovieHandler)
	mux.HandleFunc("DELETE /api/movies/{id}", movieHandler.DeleteMovieHandler)
	mux.HandleFunc("PATCH /api/movies/{id}", movieHandler.UpdateMovieHandler)

	return mux
}
