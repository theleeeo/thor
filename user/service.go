package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/theleeeo/thor/models"
	"github.com/theleeeo/thor/repo"
)

type Service struct {
	repo repo.Repo
}

func NewService(repo repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, user *models.User) (*User, error) {
	user.ID = uuid.NewString()
	user.Role = "user"

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &User{
		User: *user,
	}, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	u, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &User{
		User: *u,
	}, nil
}

func (s *Service) GetByProviderID(ctx context.Context, providerID string) (*User, error) {
	u, err := s.repo.GetUserByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}

	return &User{
		User: *u,
	}, nil
}
