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

type GenreHandler struct {
	genreService *service.GenreService
}

func NewGenreHandler(genreService *service.GenreService) *GenreHandler {
	return &GenreHandler{genreService: genreService}
}

func (h *GenreHandler) GetAllGenresHandler(w http.ResponseWriter, req *http.Request) {
	genres, err := h.genreService.ListAllGenres(req.Context())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetOneGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}

	genre, err := h.genreService.ListOneGenre(req.Context(), genreID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) CreateGenreHandler(w http.ResponseWriter, req *http.Request) {
	var genre models.Genre

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&genre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidateGenre(genre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.genreService.InsertGenre(req.Context(), &genre); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) DeleteGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil || genreID <= 0 {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}
	force := req.URL.Query().Get("force") == "true"

	if err := h.genreService.DeleteGenre(req.Context(), genreID, force); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) UpdateGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}

	var genre models.Genre
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&genre); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	genre.ID = int(genreID)

	if err := h.genreService.UpdateGenre(req.Context(), &genre); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genre)
}
