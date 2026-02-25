package services

import (
	"context"

	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
)

type MemberService interface {
	CreateMember(ctx context.Context, dB db.DB, form *dtos.Member) (*models.Member, error)
	}

type memberService struct {
	store *domain.Store
}

func NewMemberService(store *domain.Store) MemberService {
	return &memberService{store: store}
}

func (s *memberService) CreateMember(
	ctx context.Context,
	dB db.DB,
	form *dtos.Member,
) (*models.Member, error) {

	member := &models.Member{
		FullName:     form.FullName,
		IsActive:     form.IsActive,
		BenefitLimit: form.BenefitLimit,
		UsedAmount:   form.UsedAmount,
	}
	err := s.store.MemberDomain.CreateMember(ctx, dB, member)
	if err != nil {
		return nil, err
	}
	
	return member, nil
}
