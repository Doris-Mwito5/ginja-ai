package services

import (
	"context"

	"github.com/Doris-Mwito5/ginja-ai/internal/custom_types"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
)

const (
	// FraudAmountMultiplier flags a claim when requested amount exceeds procedure average cost by this factor.
	FraudAmountMultiplier = 2.0
)

type ClaimService interface {
	CreateClaim(ctx context.Context, dB db.DB, claim *dtos.ClaimSubmissionForm) (*models.Claim, error)
	GetClaimByID(ctx context.Context, dB db.DB, id int64) (*models.Claim, error)
	GetClaimByMemberID(ctx context.Context, dB db.DB, memberID string) (*models.Claim, error)
	GetClaimByProviderID(ctx context.Context, dB db.DB, providerID string) (*models.Claim, error)
	GetClaims(ctx context.Context, dB db.DB, memberID string, filter *models.Filter) (*models.ClaimList, error)
	DeleteClaim(ctx context.Context, dB db.DB, claimID int64) error
	SubmitClaim(ctx context.Context, dB db.DB, form *dtos.ClaimSubmissionForm) (*dtos.ClaimSubmissionResponse, error)
}

type claimService struct {
	store *domain.Store
}

func NewClaimService(
	store *domain.Store,
) ClaimService {
	return &claimService{
		store: store,
	}
}

func (s *claimService) SubmitClaim(
	ctx context.Context,
	dB db.DB,
	form *dtos.ClaimSubmissionForm,
) (*dtos.ClaimSubmissionResponse, error) {

	var result *dtos.ClaimSubmissionResponse
	err := dB.InTransaction(ctx, func(ctx context.Context, ops db.SQLOperations) error {
		var err error
		result, err = s.submitClaimInTx(ctx, ops, form)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *claimService) submitClaimInTx(
	ctx context.Context, 
	ops db.SQLOperations, 
	form *dtos.ClaimSubmissionForm,
	) (*dtos.ClaimSubmissionResponse, error) {

	// validate member eligibility
	member, err := s.store.MemberDomain.GetMemberByID(ctx, ops, form.MemberID)
	if err != nil || member == nil {
		return s.persistRejectedClaim(ctx, ops, form, "Member not found", false)
	}
	if !member.IsActive {
		return s.persistRejectedClaim(ctx, ops, form, "Member is not active", false)
	}

	// validate procedure exists
	procedure, err := s.store.ProcedureDomain.GetProcedureByCode(ctx, ops, form.ProcedureCode)
	if err != nil || procedure == nil {
		return s.persistRejectedClaim(ctx, ops, form, "Invalid or unknown procedure code", false)
	}

	// fraud signal: requested amount significantly above procedure average cost
	fraudFlag := form.RequestedAmount > (procedure.AverageCost * FraudAmountMultiplier)

	// check remaining benefit
	remaining := member.BenefitLimit - member.UsedAmount
	if remaining <= 0 {
		return s.persistRejectedClaim(ctx, ops, form, "Benefit limit exhausted", fraudFlag)
	}

	// determine status and approved amount
	var status custom_types.ClaimStatus
	var approvedAmount float64
	var rejectionReason string

	if form.RequestedAmount <= remaining {
		status = custom_types.ClaimStatus("APPROVED")
		approvedAmount = form.RequestedAmount
	} else {
		status = custom_types.ClaimStatus("PARTIAL")
		approvedAmount = remaining
		rejectionReason = "Requested amount exceeds remaining benefit; approved up to remaining limit."
	}

	// update member used amount
	member.UsedAmount += approvedAmount
	err = s.store.MemberDomain.CreateMember(ctx, ops, member)
	if err != nil {
		return nil, err
	}

	// persist the claim
	claim := &models.Claim{
		MemberID:        form.MemberID,
		ProviderID:      form.ProviderID,
		ProcedureCode:   form.ProcedureCode,
		DiagnosisCode:   form.DiagnosisCode,
		RequestedAmount: form.RequestedAmount,
		ApprovedAmount:  approvedAmount,
		Status:          status,
		FraudFlag:       fraudFlag,
		RejectionReason: rejectionReason,
	}
	err = s.store.ClaimDomain.CreateClaim(ctx, ops, claim)
	if err != nil {
		return nil, err
	}

	return &dtos.ClaimSubmissionResponse{
		ClaimID:         claim.ID,
		Status:          string(status),
		ApprovedAmount:  approvedAmount,
		RejectionReason: rejectionReason,
		FraudFlag:       fraudFlag,
	}, nil
}

func (s *claimService) persistRejectedClaim(
	ctx context.Context,
	ops db.SQLOperations, form *dtos.ClaimSubmissionForm,
	reason string,
	fraudFlag bool,
) (*dtos.ClaimSubmissionResponse, error) {
	claim := &models.Claim{
		MemberID:        form.MemberID,
		ProviderID:      form.ProviderID,
		ProcedureCode:   form.ProcedureCode,
		DiagnosisCode:   form.DiagnosisCode,
		RequestedAmount: form.RequestedAmount,
		ApprovedAmount:  0,
		Status:          custom_types.ClaimStatus("REJECTED"),
		FraudFlag:       fraudFlag,
		RejectionReason: reason,
	}
	if err := s.store.ClaimDomain.CreateClaim(ctx, ops, claim); err != nil {
		return nil, err
	}
	return &dtos.ClaimSubmissionResponse{
		ClaimID:         claim.ID,
		Status:          "REJECTED",
		ApprovedAmount:  0,
		RejectionReason: reason,
		FraudFlag:       fraudFlag,
	}, nil
}

func (s *claimService) CreateClaim(
	ctx context.Context, 
	dB db.DB, 
	form *dtos.ClaimSubmissionForm,
	) (*models.Claim, error) {
	claim := &models.Claim{
		MemberID:        form.MemberID,
		ProviderID:      form.ProviderID,
		ProcedureCode:   form.ProcedureCode,
		DiagnosisCode:   form.DiagnosisCode,
		RequestedAmount: form.RequestedAmount,
		ApprovedAmount:  0,
		Status:          custom_types.ClaimStatus("PENDING"),
		FraudFlag:       false,
		RejectionReason: "",
	}
	
	err := s.store.ClaimDomain.CreateClaim(ctx, dB, claim)
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func (s *claimService) GetClaimByID(
	ctx context.Context, 
	dB db.DB, 
	id int64,
	) (*models.Claim, error) {

	return s.store.ClaimDomain.GetClaimByID(ctx, dB, id)
}

func (s *claimService) GetClaimByMemberID(
	ctx context.Context, 
	dB db.DB, 
	memberID string,
	) (*models.Claim, error) {
	return s.store.ClaimDomain.GetClaimByMemberID(ctx, dB, memberID)
}

func (s *claimService) GetClaimByProviderID(
	ctx context.Context, 
	dB db.DB, 
	providerID string,
	) (*models.Claim, error) {
		
	return s.store.ClaimDomain.GetClaimByProviderID(ctx, dB, providerID)
}

func (s *claimService) GetClaims(
	ctx context.Context, 
	dB db.DB, 
	memberID string, 
	filter *models.Filter,
	) (*models.ClaimList, error) {
		
	claims, err := s.store.ClaimDomain.GetClaims(ctx, dB, memberID, filter)
	if err != nil {
		return nil, err
	}

	count, err := s.store.ClaimDomain.GetClaimsCount(ctx, dB, memberID, filter)
	if err != nil {
		return &models.ClaimList{}, err
	}
	claimList := &models.ClaimList{
		Claims: claims, 
		Pagination: models.NewPagination(count, filter.Page, filter.Per),
	}

	return claimList, nil
}

func (s *claimService) DeleteClaim(ctx context.Context, dB db.DB, claimID int64) error {
	return s.store.ClaimDomain.DeleteClaim(ctx, dB, claimID)
}
