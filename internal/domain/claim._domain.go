package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
	"github.com/Doris-Mwito5/ginja-ai/internal/null"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
)

const (
	createClaimSQL          = "INSERT INTO claims (member_id, provider_id, procedure_code, diagnosis_code, requested_amount, approved_amount, status, fraud_flag, rejection_reason) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id"
	getClaimsSQL            = "SELECT id, member_id, provider_id, procedure_code, diagnosis_code, requested_amount, approved_amount, status, fraud_flag, rejection_reason, created_at, updated_at FROM claims"
	getClaimByIDSQL         = getClaimsSQL + " WHERE id = $1"
	getClaimByMemberIDSQL   = getClaimsSQL + " WHERE member_id = $1"
	getClaimByProviderIDSQL = getClaimsSQL + " WHERE provider_id = $1"
	getClaimsCountSQL       = "SELECT COUNT(*) FROM claims"
	updateClaimSQL          = "UPDATE claims SET member_id = $1, provider_id = $2, procedure_code = $3, diagnosis_code = $4, requested_amount = $5, approved_amount = $6, status = $7, fraud_flag = $8, rejection_reason = $9 WHERE id = $10"
	deleteClaimSQL          = "DELETE FROM claims WHERE id = $1"
)

type (
	ClaimDomain interface {
		CreateClaim(ctx context.Context, operations db.SQLOperations, claim *models.Claim) error
		GetClaimByID(ctx context.Context, operations db.SQLOperations, id int64) (*models.Claim, error)
		GetClaimByMemberID(ctx context.Context, operations db.SQLOperations, memberID string) (*models.Claim, error)
		GetClaimByProviderID(ctx context.Context, operations db.SQLOperations, providerID string) (*models.Claim, error)
		GetClaimsCount(ctx context.Context, operations db.SQLOperations, memberID string, filter *models.Filter) (int, error)
		GetClaims(ctx context.Context, operations db.SQLOperations, memberID string, filter *models.Filter) ([]*models.Claim, error)
		DeleteClaim(ctx context.Context, operations db.SQLOperations, claimID int64) error
	}

	claimDomain struct{}
)

func NewClaimDomain() ClaimDomain {
	return &claimDomain{}
}

func (s *claimDomain) CreateClaim(
	ctx context.Context,
	operations db.SQLOperations,
	claim *models.Claim) error {

	claim.Touch()
	if claim.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createClaimSQL,
			claim.MemberID,
			claim.ProviderID,
			claim.ProcedureCode,
			claim.DiagnosisCode,
			claim.RequestedAmount,
			claim.ApprovedAmount,
			claim.Status,
			claim.FraudFlag,
			claim.RejectionReason,
		).Scan(&claim.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("create claim query error: %v", err)
		}
		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateClaimSQL,
		claim.MemberID,
		claim.ProviderID,
		claim.ProcedureCode,
		claim.DiagnosisCode,
		claim.RequestedAmount,
		claim.ApprovedAmount,
		claim.Status,
		claim.FraudFlag,
		claim.RejectionReason,
		claim.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update claim query error: %v", err)
	}

	return nil
}

func (s *claimDomain) GetClaimByID(
	ctx context.Context,
	operations db.SQLOperations,
	id int64) (*models.Claim, error) {

	row := operations.QueryRowContext(
		ctx,
		getClaimByIDSQL,
		id,
	)

	return s.scanRow(row)
}

func (s *claimDomain) GetClaimByProviderID(
	ctx context.Context,
	operations db.SQLOperations,
	providerID string) (*models.Claim, error) {

	row := operations.QueryRowContext(
		ctx,
		getClaimByProviderIDSQL,
		providerID,
	)

	return s.scanRow(row)
}

func (s *claimDomain) GetClaimByMemberID(
	ctx context.Context,
	operations db.SQLOperations,
	memberID string) (*models.Claim, error) {

	row := operations.QueryRowContext(
		ctx,
		getClaimByMemberIDSQL,
		memberID,
	)
	return s.scanRow(row)
}

func (s *claimDomain) GetClaimsCount(
	ctx context.Context,
	operations db.SQLOperations,
	memberID string,
	filter *models.Filter,
) (int, error) {

	filter.MemberID = null.NullValue(memberID)
	query, args := s.buildQuery(getClaimsCountSQL, filter.NoPagination())

	rows := operations.QueryRowContext(
		ctx,
		query,
		args...,
	)

	var count int

	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(
			err,
		).LogErrorMessage(
			"claims count by query row contex err: %v",
			err,
		)
	}
	return count, nil

}

func (s *claimDomain) GetClaims(
	ctx context.Context,
	operations db.SQLOperations,
	memberID string,
	filter *models.Filter,
) ([]*models.Claim, error) {

	filter.MemberID = null.NullValue(memberID)
	query, args := s.buildQuery(getClaimsSQL, filter.NoPagination())
	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return []*models.Claim{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("claims query row err: %v", err)
	}

	defer rows.Close()

	claims := make([]*models.Claim, 0)

	for rows.Next() {
		claim, err := s.scanRow(rows)
		if err != nil {
			return []*models.Claim{}, err
		}
		claims = append(claims, claim)
	}
	if rows.Err() != nil {
		return []*models.Claim{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage(
			"list claims err: %v",
			err,
		)
	}

	return claims, nil
}

func (s *claimDomain) DeleteClaim(
	ctx context.Context,
	operations db.SQLOperations,
	claimID int64,
) error {

	_, err := operations.ExecContext(
		ctx,
		deleteClaimSQL,
		claimID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("delete claim err")
	}

	return nil
}

func (s *claimDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {

	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.MemberID != nil {
		condition := fmt.Sprintf("member_id = $%d", counter.Touch())
		args = append(args, null.ValueFromNull(filter.MemberID))
		conditions = append(conditions, condition)
	}

	if filter.ProviderID != nil {
		condition := fmt.Sprintf("provider_id = $%d", counter.Touch())
		args = append(args, null.ValueFromNull(filter.ProviderID))
		conditions = append(conditions, condition)
	}

	if filter.Term != "" {
		textCols := []string{"procedure_code", "diagnosis_code", "rejection_reason"}
		likeStatements := make([]string, 0)
		term := strings.ToLower(filter.Term)

		for _, col := range textCols {
			likeStmt := fmt.Sprintf(
				" (LOWER(%s) LIKE '%%' || $%d || '%%') ", col, counter.Touch())
			likeStatements = append(likeStatements, likeStmt)
			args = append(args, term)
		}

		conditions = append(conditions, " ("+strings.Join(likeStatements, " OR ")+")")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	if filter.Page > 0 && filter.Per > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", counter.Touch(), counter.Touch())
		args = append(args, filter.Per, (filter.Page-1)*filter.Per)
	}

	return query, args
}

func (s *claimDomain) scanRow(
	row db.RowScanner,
) (*models.Claim, error) {

	var claim models.Claim
	err := row.Scan(
		&claim.ID,
		&claim.MemberID,
		&claim.ProviderID,
		&claim.ProcedureCode,
		&claim.DiagnosisCode,
		&claim.RequestedAmount,
		&claim.ApprovedAmount,
		&claim.Status,
		&claim.FraudFlag,
		&claim.RejectionReason,
		&claim.CreatedAt,
		&claim.UpdatedAt,
	)
	if err != nil {
		return &models.Claim{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan row error: %v", err)
	}
	return &claim, nil
}
