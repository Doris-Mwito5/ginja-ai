package services

import (
	"context"

	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
)

type (
	ProcedureService interface {
		CreateProcedure(ctx context.Context, dB db.DB, form *dtos.Procedure) (*models.Procedure, error)
	}

	procedureService struct {
		store *domain.Store
	}
) 


func NewProcedureService(store *domain.Store) ProcedureService {
	return &procedureService{store: store}
}

func (s *procedureService) CreateProcedure(
	ctx context.Context, 
	dB db.DB, 
	form *dtos.Procedure,
) (*models.Procedure, error) {

	procedure := &models.Procedure{
		Code:        form.Code,
		Description: form.Description,
		AverageCost: form.AverageCost,
	}
	err := s.store.ProcedureDomain.CreateProcedure(ctx, dB, procedure)
	if err != nil {
		return nil, err
	}

	return procedure, nil
}

