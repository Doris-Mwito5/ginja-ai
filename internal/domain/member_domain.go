package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/models"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
)

const (
	createMemberSQL        = "INSERT INTO members (full_name, is_active, benefit_limit, used_amount) VALUES ($1, $2, $3, $4) RETURNING id"
	getMembersSQL          = "SELECT id, full_name, is_active, benefit_limit, used_amount, created_at, updated_at FROM members"
	getMemberByIDSQL       = getMembersSQL + " WHERE id = $1"
	getMemberByFullNameSQL = getMembersSQL + " WHERE full_name = $1"
	getMembersCountSQL     = "SELECT COUNT(*) FROM members"
	updateMemberSQL        = "UPDATE members SET full_name = $1, is_active = $2, benefit_limit = $3, used_amount = $4 WHERE id = $5"
	deleteMemberSQL        = "DELETE FROM members WHERE id = $1"
)

type (
	MemberDomain interface {
		CreateMember(ctx context.Context, operations db.SQLOperations, member *models.Member) error
		GetMemberByID(ctx context.Context, operations db.SQLOperations, id int64) (*models.Member, error)
		GetMemberByFullName(ctx context.Context, operations db.SQLOperations, fullName string) (*models.Member, error)
		GetMembersCount(ctx context.Context, operations db.SQLOperations, filter *models.Filter) (int, error)
		GetMembers(ctx context.Context, operations db.SQLOperations, filter *models.Filter) ([]*models.Member, error)
		DeleteMember(ctx context.Context, operations db.SQLOperations, id int64) error
	}

	memberDomain struct{}
)

func NewMemberDomain() MemberDomain {
	return &memberDomain{}
}

func (s *memberDomain) CreateMember(
	ctx context.Context,
	operations db.SQLOperations,
	member *models.Member,
) error {
	member.Touch()

	if member.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createMemberSQL,
			member.FullName,
			member.IsActive,
			member.BenefitLimit,
			member.UsedAmount,
		).Scan(&member.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("create member query error: %v", err)
		}
		return nil
	}

	_, err := operations.ExecContext(
		ctx,
		updateMemberSQL,
		member.FullName,
		member.IsActive,
		member.BenefitLimit,
		member.UsedAmount,
		member.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update member query error: %v", err)
	}

	return nil
}

func (s *memberDomain) GetMemberByID(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) (*models.Member, error) {
	row := operations.QueryRowContext(
		ctx,
		getMemberByIDSQL,
		id,
	)

	return s.scanRow(row)
}

func (s *memberDomain) GetMemberByFullName(
	ctx context.Context,
	operations db.SQLOperations,
	fullName string,
) (*models.Member, error) {
	row := operations.QueryRowContext(
		ctx,
		getMemberByFullNameSQL,
		fullName,
	)

	return s.scanRow(row)
}

func (s *memberDomain) GetMembersCount(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) (int, error) {
	query, args := s.buildQuery(getMembersCountSQL, filter.NoPagination())
	rows := operations.QueryRowContext(ctx, query, args...)
	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get members count query error: %v", err)
	}
	return count, nil
}

func (s *memberDomain) GetMembers(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) ([]*models.Member, error) {

	query, args := s.buildQuery(getMembersSQL, filter.NoPagination())

	rows, err := operations.QueryContext(
		ctx, 
		query, 
		args...,
	)
	if err != nil {
		return []*models.Member{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get members query error: %v", err)
	}
	
	defer rows.Close()

	members := make([]*models.Member, 0)

	for rows.Next() {
		member, err := s.scanRow(rows)
		if err != nil {
			return []*models.Member{}, err
		}
		members = append(members, member)
	}
	
	if rows.Err() != nil {
		return []*models.Member{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage("list members err: %v", err)
	}

	return members, nil
}

func (s *memberDomain) DeleteMember(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) error {
	_, err := operations.ExecContext(
		ctx, 
		deleteMemberSQL, 
		id,
	)

	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("delete member query error: %v", err)
	}
	return nil
}

func (s *memberDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {

	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		textCols := []string{"full_name"}
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

func (s *memberDomain) scanRow(
	row db.RowScanner,
) (*models.Member, error) {
	
	var member models.Member
	err := row.Scan(
		&member.ID,
		&member.FullName,
		&member.IsActive,
		&member.BenefitLimit,
		&member.UsedAmount,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return &models.Member{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan row error: %v", err)
	}
	return &member, nil
}
