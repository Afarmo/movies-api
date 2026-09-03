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

	var movies []*models.Movie
	var filter models.MovieFilter

	if genreStr := query.Get("genre"); genreStr != "" {
		genreID, err := strconv.ParseInt(genreStr, 10, 64)
		if err != nil || genreID < 1 {
			http.Error(w, "invalid genre id", http.StatusBadRequest)
			return
		}

		filter.GenreID = &genreID
	}

	if actorStr := query.Get("actor"); actorStr != "" {
		actorID, err := strconv.ParseInt(actorStr, 10, 64)
		if err != nil || actorID < 1 {
			http.Error(w, "invalid actor id", http.StatusBadRequest)
			return
		}

		filter.ActorID = &actorID
	}

	if yearStr := query.Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "invalid release year", http.StatusBadRequest)
			return
		}

		filter.Year = &year
	}

	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	if pageStr != "" && sizeStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "invalid size", http.StatusBadRequest)
			return
		}

		filter.Page = &page
		filter.Size = &size
	}

	movies, err := h.movieService.ListMovies(req.Context(), &filter)
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

	var movies []*models.Movie
	var filter models.MovieFilter
	var err error

	if title := req.URL.Query().Get("title"); title != "" {
		movies, err = h.movieService.SearchMovies(req.Context(), title)
	} else {
		movies, err = h.movieService.ListMovies(req.Context(), &filter)
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
