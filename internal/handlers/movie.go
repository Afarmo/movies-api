package handlers

import (
	"encoding/json"
	"fmt"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (h *MovieHandler) GetAllMoviesHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	genreStr := query.Get("genre")
	yearStr := query.Get("year")
	actorStr := query.Get("actor")

	var movies []*models.Movie
	var err error

	switch {
	case genreStr != "":
		var genreID int64

		genreID, err = strconv.ParseInt(genreStr, 10, 64)
		if err == nil {
			movies, err = h.movieService.MoviesByGenre(req.Context(), genreID)
		}

	case yearStr != "":
		var releaseYear int

		releaseYear, err = strconv.Atoi(yearStr)
		if err == nil {
			movies, err = h.movieService.MoviesByYear(req.Context(), releaseYear)
		}

	case actorStr != "":
		var actorID int64

		actorID, err = strconv.ParseInt(actorStr, 10, 64)
		if err == nil {
			movies, err = h.movieService.MoviesByActor(req.Context(), actorID)
		}

	default:
		movies, err = h.movieService.ListAllMovies(req.Context())
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) GetOneMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	movie, err := h.movieService.ListOneMovie(req.Context(), movieID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) CreateMovieHandler(w http.ResponseWriter, req *http.Request) {
	var movie models.Movie

	err := json.NewDecoder(req.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateMovie(movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(movie.ActorIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(movie.GenreIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.movieService.InsertMovie(req.Context(), &movie); err != nil {
		fmt.Println(movie)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	if err := h.movieService.DeleteMovie(req.Context(), movieID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil || movieID <= 0 {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	var movie models.MoviePatch

	if err := json.NewDecoder(req.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := validation.ValidateMoviePatch(movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(movie.AddActorIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(movie.AddGenreIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidateIDs(movie.RemoveActorIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(movie.RemoveGenreIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.movieService.UpdateMovie(req.Context(), movieID, &movie); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&movie)
}

func (h *MovieHandler) SearchMoviesHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")

	var movies []*models.Movie
	var err error

	switch {
	case title != "":
		movies, err = h.movieService.SearchMovies(req.Context(), title)
	default:
		movies, err = h.movieService.ListAllMovies(req.Context())
	}

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
