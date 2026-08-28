package handlers

import (
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) GetOneGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid genre id", http.StatusBadRequest)
		return
	}

	genre, err := h.genreService.ListOneGenre(req.Context(), genreID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) CreateGenreHandler(w http.ResponseWriter, req *http.Request) {
	var genre models.Genre

	if err := json.NewDecoder(req.Body).Decode(&genre); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.genreService.InsertGenre(req.Context(), &genre); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) DeleteGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid genre id", http.StatusBadRequest)
		return
	}

	if err := h.genreService.DeleteGenre(req.Context(), genreID); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) UpdateGenreHandler(w http.ResponseWriter, req *http.Request) {
	genreID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid genre id", http.StatusBadRequest)
		return
	}

	var genre models.Genre

	if err := json.NewDecoder(req.Body).Decode(&genre); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	genre.ID = int(genreID)

	if err := h.genreService.UpdateGenre(req.Context(), &genre); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
