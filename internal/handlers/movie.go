package handlers

import (
	"encoding/json"
	"movies-api/internal/apperrors"
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

	var filter models.MovieFilter

	if genreStr := query.Get("genre"); genreStr != "" {
		genreID, err := strconv.ParseInt(genreStr, 10, 64)
		if err != nil || genreID < 1 {
			http.Error(w, "Invalid genre ID", http.StatusBadRequest)
			return
		}

		filter.GenreID = &genreID
	}

	if actorStr := query.Get("actor"); actorStr != "" {
		actorID, err := strconv.ParseInt(actorStr, 10, 64)
		if err != nil || actorID < 1 {
			http.Error(w, "Invalid actor ID", http.StatusBadRequest)
			return
		}

		filter.ActorID = &actorID
	}

	if yearStr := query.Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			http.Error(w, "Invalid release year", http.StatusBadRequest)
			return
		}

		filter.Year = &year
	}

	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	if pageStr != "" || sizeStr != "" {
		if pageStr == "" || sizeStr == "" {
			http.Error(w, "page and size must be provided together", http.StatusBadRequest)
			return
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			http.Error(w, "Invalid page", http.StatusBadRequest)
			return
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}

		filter.Page = &page
		filter.Size = &size
	}

	movies, err := h.movieService.ListMovies(req.Context(), &filter)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	totalPages := 0

	if filter.Page != nil && filter.Size != nil {
		totalPages = (movies.Total + *filter.Size - 1) / *filter.Size
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ListMoviesResponse{
		Movies:     movies.Movies,
		Page:       filter.Page,
		Size:       filter.Size,
		Total:      movies.Total,
		TotalPages: totalPages,
	})
}

func (h *MovieHandler) GetOneMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}
	movie, err := h.movieService.ListOneMovie(req.Context(), movieID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) CreateMovieHandler(w http.ResponseWriter, req *http.Request) {
	var movie models.Movie

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) DeleteMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	if err := h.movieService.DeleteMovie(req.Context(), movieID); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) UpdateMovieHandler(w http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil || movieID <= 0 {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	var moviePatch models.MoviePatch

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&moviePatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidateMoviePatch(moviePatch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(moviePatch.AddActorIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(moviePatch.AddGenreIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidateIDs(moviePatch.RemoveActorIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validation.ValidateIDs(moviePatch.RemoveGenreIds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movie, err := h.movieService.UpdateMovie(req.Context(), movieID, &moviePatch)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&movie)
}

func (h *MovieHandler) SearchMoviesHandler(w http.ResponseWriter, req *http.Request) {
	title := req.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	movies, err := h.movieService.SearchMovies(req.Context(), title)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
