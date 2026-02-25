package services

import (
	"context"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
)

type UserService interface {
	Register(ctx context.Context, dB db.DB, form *dtos.RegisterRequest) (*models.User, error)
	GetUserByID(ctx context.Context, dB db.DB, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, dB db.DB, username string) (*models.User, error)
	ValidateCredentials(ctx context.Context, dB db.DB, username, password string) (*models.User, error)
}

type userService struct {
	store *domain.Store
}

func NewUserService(store *domain.Store) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) Register(
	ctx context.Context, 
	dB db.DB, form *dtos.RegisterRequest,
	) (*models.User, error) {
	
	if form.Username == "" || form.Email == "" {
		return nil, apperr.NewBadRequest("username and email are required")
	}

	existing, _ := s.store.UserDomain.GetUserByUsername(ctx, dB, form.Username)
	if existing != nil {
		return nil, apperr.NewConflict("username", form.Username)
	}

	existing, _ = s.store.UserDomain.GetUserByEmail(ctx, dB, form.Email)
	if existing != nil {
		return nil, apperr.NewConflict("email", form.Email)
	}

	passwordHash, err := utils.HashPassword(form.Password)
	if err != nil {
		return nil, apperr.NewInternal("failed to process password")
	}

	user := &models.User{
		Username:     form.Username,
		Email:        form.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
	}
	if err := s.store.UserDomain.CreateUser(ctx, dB, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(
	ctx context.Context, 
	dB db.DB, 
	id int64,
) (*models.User, error) {

	user, err := s.store.UserDomain.GetUserByID(ctx, dB, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByUsername(
	ctx context.Context,
	dB db.DB,
	username string,
) (*models.User, error) {

	user, err := s.store.UserDomain.GetUserByUsername(ctx, dB, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ValidateCredentials(
	ctx context.Context,
	dB db.DB,
	username,
	password string,
) (*models.User, error) {

	user, err := s.store.UserDomain.GetUserByUsername(ctx, dB, username)
	if err != nil || user == nil {
		return nil, apperr.NewAuthorization("invalid username or password")
	}

	if !user.IsActive {
		return nil, apperr.NewAuthorization("account is inactive")
	}

	err = utils.ValidatePassword(password, user.PasswordHash)
	if err != nil {
		return nil, apperr.NewAuthorization("invalid username or password")
	}
	return user, nil
}
