package service

import (
	"context"
	"movies-api/internal/models"
	"movies-api/internal/repository"
)

type ActorService struct {
	repo *repository.ActorRepo
}

func NewActorService(repo *repository.ActorRepo) *ActorService {
	return &ActorService{repo: repo}
}

func (s *ActorService) InsertActor(ctx context.Context, actor *models.Actor) error {
	return s.repo.InsertActor(ctx, actor)
}

func (s *ActorService) ListAllActors(ctx context.Context) ([]*models.Actor, error) {
	return s.repo.ListAllActors(ctx)
}

func (s *ActorService) ListOneActor(ctx context.Context, id int64) (*models.Actor, error) {
	return s.repo.ListOneActor(ctx, id)
}

func (s *ActorService) UpdateActor(ctx context.Context, actor *models.Actor) error {
	return s.repo.UpdateActor(ctx, actor)
}

func (s *ActorService) DeleteActor(ctx context.Context, id int64) error {
	return s.repo.DeleteActor(ctx, id)
}

func (s *ActorService) CountActors(ctx context.Context) (int64, error) {
	return s.repo.CountActors(ctx)
}

func (s *ActorService) SearchActors(ctx context.Context, name string) ([]*models.Actor, error) {
	return s.repo.SearchActors(ctx, name)
}
