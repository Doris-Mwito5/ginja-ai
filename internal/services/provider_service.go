package services

import (
	"context"

	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
)

type (
	ProviderService interface {
		CreateProvider(ctx context.Context, dB db.DB, form *dtos.Provider) (*models.Provider, error)
	}

	providerService struct {
		store *domain.Store
	}
)

func NewProviderService(store *domain.Store) ProviderService {
	return &providerService{store: store}
}
func (s *providerService) CreateProvider(
ctx context.Context, 
dB db.DB, 
form *dtos.Provider,
) (*models.Provider, error) {

	provider := &models.Provider{
		Name:     form.Name,
		Location: form.Location,
	}

	err := s.store.ProviderDomain.CreateProvider(ctx, dB, provider)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

