package handlers

import (
	"movies-api/internal/service"
	"net/http"
)

type GenreHandler struct {
	genreService *service.GenreService
}

func NewGenreHandler(genreService *service.GenreService) *GenreHandler {
	return &GenreHandler{genreService: genreService}
}

func (h *GenreHandler) GetOneGenreHandler(w http.ResponseWriter, req *http.Request) {

}
