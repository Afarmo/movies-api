package service

import (
	"context"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type GenreService struct {
	repo *repository.GenreRepo
}

func NewGenreService(repo *repository.GenreRepo) *GenreService {
	return &GenreService{repo: repo}
}

func (s *GenreService) InsertGenre(ctx context.Context, genre *models.Genre) error {
	return s.repo.InsertGenre(ctx, genre)
}

func (s *GenreService) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	return s.repo.UpdateGenre(ctx, genre)
}

func (s *GenreService) DeleteGenre(ctx context.Context, id int64) error {
	return s.repo.DeleteGenre(ctx, id)
}

func (s *GenreService) ListAllGenres(ctx context.Context) ([]*models.Genre, error) {
	return s.repo.ListAllGenres(ctx)
}

func (s *GenreService) ListOneGenre(ctx context.Context, id int64) (*models.Genre, error) {
	return s.repo.ListOneGenre(ctx, id)
}
