package service

import (
	"context"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type MovieService struct {
	repo *repository.MovieRepo
}

func NewMovieService(repo *repository.MovieRepo) *MovieService {
	return &MovieService{repo: repo}
}

func (s *MovieService) InsertMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.InsertMovie(ctx, movie)
}

func (s *MovieService) ListAllMovies(ctx context.Context) ([]*models.Movie, error) {
	return s.repo.ListAllMovies(ctx)
}

func (s *MovieService) ListOneMovie(ctx context.Context, id int64) (*models.Movie, error) {
	return s.repo.ListOneMovie(ctx, id)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.UpdateMovie(ctx, movie)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int64) error {
	return s.repo.DeleteMovie(ctx, id)
}

func (s *MovieService) SearchMovies(ctx context.Context, title string) ([]*models.Movie, error) {
	return s.repo.SearchMovies(ctx, title)
}
