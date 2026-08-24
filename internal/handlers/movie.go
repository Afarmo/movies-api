package handlers

import (
	"database/sql"
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type MovieHandler struct {
	movieService *service.MovieService
}

func NewMovieHandler(genreService *service.MovieService) *MovieHandler {
	return &MovieHandler{movieService: genreService}
}

func (h *MovieHandler) GetOneMovieHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idString := req.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "invalid movie id", http.StatusBadRequest)
		return
	}
	movie, err := h.movieService.ListOneMovie(ctx, id)
	if err != nil {
		return // TODO
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) GetAllMoviesHAndler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	movies, err := h.movieService.ListAllMovies(ctx)
	if err != nil {
		return // TODO
	}
	w.Header().Set("Conent-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func (h *MovieHandler) CreateMovieHandler(w http.ResponseWriter, req *http.Request) {
	movie := models.Movie{}
	err := json.NewDecoder(req.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "could not create movie", http.StatusInternalServerError)
		return
	}
	ctx := req.Context()
	err = h.movieService.InsertMovie(ctx, &movie)
	if err != nil {
		http.Error(w, "could not create movie", http.StatusInternalServerError)
		return // TODO
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	ctx := req.Context()

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.movieService.DeleteMovie(ctx, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idString := req.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	movie := models.Movie{}
	movie.ID = int(id)
	err = json.NewDecoder(req.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	err = h.movieService.UpdateMovie(ctx, &movie)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "No such movie", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *MovieHandler) GetMovieByNameHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	name := req.PathValue("name")
	movie, err := h.movieService.MovieByName(ctx, name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(movie)
}
