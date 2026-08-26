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

func (h *GenreHandler) GetOneGenreHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	idString := req.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return // TODO
	}

	genre, err := h.genreService.ListOneGenre(ctx, id)
	if err != nil {
		return // TODO
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) GetAllGenresHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	genres, err := h.genreService.ListAllGenres(ctx)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genres)
}

func (h *GenreHandler) CreateGenreHandler(w http.ResponseWriter, req *http.Request) {

	genre := models.Genre{}
	err := json.NewDecoder(req.Body).Decode(&genre)
	if err != nil {
		return // TODO
	}

	ctx := req.Context()
	err = h.genreService.InsertGenre(ctx, &genre)
	if err != nil {
		return // TODO
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(genre)
}

func (h *GenreHandler) DeleteGenreHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	idString := req.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return // TODO
	}
	err = h.genreService.DeleteGenre(ctx, id)
	if err != nil {
		return // TODO
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *GenreHandler) UpdateGenreHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idString := req.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return // TODO
	}

	genre := models.Genre{}
	genre.ID = int(id)
	err = json.NewDecoder(req.Body).Decode(&genre)
	if err != nil {
		return // TOOD
	}
	err = h.genreService.UpdateGenre(ctx, &genre)
	if err != nil {
		return // TODO
	}
	w.WriteHeader(http.StatusAccepted)
}
