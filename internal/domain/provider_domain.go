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
	createProviderSQL    = "INSERT INTO providers (name, location) VALUES ($1, $2) RETURNING id"
	getProvidersSQL      = "SELECT id, name, location, created_at, updated_at FROM providers"
	getProviderByIDSQL   = getProvidersSQL + " WHERE id = $1"
	getProviderByNameSQL = getProvidersSQL + " WHERE name = $1"
	getProvidersCountSQL = "SELECT COUNT(*) FROM providers"
	updateProviderSQL    = "UPDATE providers SET name = $1, location = $2 WHERE id = $3"
	deleteProviderSQL    = "DELETE FROM providers WHERE id = $1"
)

type (
	ProviderDomain interface {
		CreateProvider(ctx context.Context, operations db.SQLOperations, provider *models.Provider) error
		GetProviderByID(ctx context.Context, operations db.SQLOperations, id int64) (*models.Provider, error)
		GetProviderByName(ctx context.Context, operations db.SQLOperations, name string) (*models.Provider, error)
		GetProvidersCount(ctx context.Context, operations db.SQLOperations, filter *models.Filter) (int, error)
		GetProviders(ctx context.Context, operations db.SQLOperations, filter *models.Filter) ([]*models.Provider, error)
		DeleteProvider(ctx context.Context, operations db.SQLOperations, id int64) error
	}

	providerDomain struct{}
)

func NewProviderDomain() ProviderDomain {
	return &providerDomain{}
}

func (s *providerDomain) CreateProvider(
	ctx context.Context,
	operations db.SQLOperations,
	provider *models.Provider,
) error {
	provider.Touch()

	if provider.IsNew() {
		err := operations.QueryRowContext(
			ctx,
			createProviderSQL,
			provider.Name,
			provider.Location,
		).Scan(&provider.ID)
		if err != nil {
			return apperr.NewDatabaseError(
				err,
			).LogErrorMessage("create provider query error: %v", err)
		}
		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateProviderSQL,
		provider.Name,
		provider.Location,
		provider.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(
			err,
		).LogErrorMessage("update provider query error: %v", err)
	}
	return nil
}

func (s *providerDomain) GetProviderByID(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) (*models.Provider, error) {
	row := operations.QueryRowContext(
		ctx,
		getProviderByIDSQL,
		id,
	)
	return s.scanRow(row)
}

func (s *providerDomain) GetProviderByName(
	ctx context.Context,
	operations db.SQLOperations,
	name string,
) (*models.Provider, error) {

	row := operations.QueryRowContext(
		ctx,
		getProviderByNameSQL,
		name,
	)

	return s.scanRow(row)
}	

func (s *providerDomain) GetProvidersCount(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) (int, error) {
	query, args := s.buildQuery(getProvidersCountSQL, filter.NoPagination())
	rows := operations.QueryRowContext(
		ctx, 
		query, 
		args...,
	)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(err).LogErrorMessage("get providers count query error: %v", err)
	}

	return count, nil
}

func (s *providerDomain) GetProviders(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) ([]*models.Provider, error) {
	query, args := s.buildQuery(getProvidersSQL, filter.NoPagination())
	rows, err := operations.QueryContext(
		ctx, 
		query, 
		args...,
	)
	if err != nil {
		return []*models.Provider{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("get providers query error: %v", err)
	}
	defer rows.Close()

	providers := make([]*models.Provider, 0)
	for rows.Next() {
		provider, err := s.scanRow(rows)
		if err != nil {
			return []*models.Provider{}, err
		}
		providers = append(providers, provider)
	}

	if rows.Err() != nil {
		return []*models.Provider{}, apperr.NewDatabaseError(
			rows.Err(),
		).LogErrorMessage("list providers err: %v", err)
	}
	return providers, nil
}

func (s *providerDomain) DeleteProvider(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) error {
	_, err := operations.ExecContext(
		ctx,
		deleteProviderSQL,
		id,
	)
	if err != nil {
		return apperr.NewDatabaseError(err).LogErrorMessage("delete provider query error: %v", err)
	}
	return nil
}

func (s *providerDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {
	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		textCols := []string{"name", "location"}
		likeStatements := make([]string, 0)
		term := strings.ToLower(filter.Term)
		for _, col := range textCols {
			likeStmt := fmt.Sprintf(" (LOWER(%s) LIKE '%%' || $%d || '%%') ", col, counter.Touch())
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

func (s *providerDomain) scanRow(
	row db.RowScanner,
) (*models.Provider, error) {
	
	var provider models.Provider
	err := row.Scan(
		&provider.ID,
		&provider.Name,
		&provider.Location,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err != nil {
		return &models.Provider{}, apperr.NewDatabaseError(
			err,
		).LogErrorMessage("scan row error: %v", err)
	}
	return &provider, nil
}