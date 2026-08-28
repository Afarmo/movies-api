package handlers

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
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
	movies, err := h.movieService.ListAllMovies(req.Context())
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
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.movieService.InsertMovie(req.Context(), &movie); err != nil {
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
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}

	var movie models.MoviePatch

	if err := json.NewDecoder(req.Body).Decode(&movie); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.movieService.UpdateMovie(req.Context(), movieID, &movie); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *MovieHandler) SearchMoviesHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")

	movies, err := h.movieService.SearchMovies(req.Context(), title)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
