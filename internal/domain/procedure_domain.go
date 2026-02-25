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
	createProcedureSQL    = "INSERT INTO procedures (code, description, average_cost) VALUES ($1, $2, $3) RETURNING id"
	getProceduresSQL      = "SELECT id, code, description, average_cost, created_at, updated_at FROM procedures"
	getProcedureByIDSQL   = getProceduresSQL + " WHERE id = $1"
	getProcedureByCodeSQL = getProceduresSQL + " WHERE code = $1"
	getProceduresCountSQL = "SELECT COUNT(*) FROM procedures"
	updateProcedureSQL    = "UPDATE procedures SET code = $1, description = $2, average_cost = $3 WHERE id = $4"
	deleteProcedureSQL    = "DELETE FROM procedures WHERE id = $1"
)

type (
	ProcedureDomain interface {
		CreateProcedure(ctx context.Context, operations db.SQLOperations, procedure *models.Procedure) error
		GetProcedureByID(ctx context.Context, operations db.SQLOperations, id int64) (*models.Procedure, error)
		GetProcedureByCode(ctx context.Context, operations db.SQLOperations, code string) (*models.Procedure, error)
		GetProceduresCount(ctx context.Context, operations db.SQLOperations, filter *models.Filter) (int, error)
		GetProcedures(ctx context.Context, operations db.SQLOperations, filter *models.Filter) ([]*models.Procedure, error)
		DeleteProcedure(ctx context.Context, operations db.SQLOperations, id int64) error
	}

	procedureDomain struct{}
)

func NewProcedureDomain() ProcedureDomain {
	return &procedureDomain{}
}

func (s *procedureDomain) CreateProcedure(
	ctx context.Context,
	operations db.SQLOperations,
	procedure *models.Procedure,
) error {

	procedure.Touch()

	if procedure.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createProcedureSQL,
			procedure.Code,
			procedure.Description,
			procedure.AverageCost,
		).Scan(&procedure.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("create procedure query error: %v", err)
		}
		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateProcedureSQL,
		procedure.Code,
		procedure.Description,
		procedure.AverageCost,
		procedure.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update procedure query error: %v", err)
	}
	return nil
}

func (s *procedureDomain) GetProcedureByID(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) (*models.Procedure, error) {
	row := operations.QueryRowContext(
		ctx,
		getProcedureByIDSQL,
		id,
	)
	return s.scanRow(row)
}

func (s *procedureDomain) GetProcedureByCode(
	ctx context.Context,
	operations db.SQLOperations,
	code string,
) (*models.Procedure, error) {
	row := operations.QueryRowContext(
		ctx,
		getProcedureByCodeSQL,
		code,
	)
	return s.scanRow(row)
}

func (s *procedureDomain) GetProceduresCount(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) (int, error) {

	query, args := s.buildQuery(getProceduresCountSQL, filter.NoPagination())

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
		).LogErrorMessage("get procedures count query error: %v", err)
	}
	return count, nil
}

func (s *procedureDomain) GetProcedures(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) ([]*models.Procedure, error) {

	query, args := s.buildQuery(getProceduresSQL, filter.NoPagination())

	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return []*models.Procedure{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get procedures query error: %v", err)
	}
	defer rows.Close()

	procedures := make([]*models.Procedure, 0)

	for rows.Next() {
		procedure, err := s.scanRow(rows)
		if err != nil {
			return []*models.Procedure{}, err
		}
		procedures = append(procedures, procedure)
	}
	if rows.Err() != nil {
		return []*models.Procedure{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage(
			"list procedures err: %v",
			err,
		)
	}

	return procedures, nil
}

func (s *procedureDomain) DeleteProcedure(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) error {
	_, err := operations.ExecContext(
		ctx,
		deleteProcedureSQL,
		id,
	)

	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("delete procedure query error: %v", err)
	}
	return nil
}

func (s *procedureDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {
	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		textCols := []string{"code", "description"}
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

	return query, args
}

func (s *procedureDomain) scanRow(
	row db.RowScanner,
) (*models.Procedure, error) {
	var procedure models.Procedure
	err := row.Scan(
		&procedure.ID,
		&procedure.Code,
		&procedure.Description,
		&procedure.AverageCost,
		&procedure.CreatedAt,
		&procedure.UpdatedAt,
	)
	if err != nil {
		return &models.Procedure{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan row error: %v", err)
	}
	return &procedure, nil
}
