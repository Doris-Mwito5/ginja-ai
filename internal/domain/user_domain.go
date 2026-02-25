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
	createUserSQL        = "INSERT INTO users (username, email, password_hash, is_active) VALUES ($1, $2, $3, $4) RETURNING id"
	getUsersSQL          = "SELECT id, username, email, password_hash, is_active, created_at, updated_at FROM users"
	getUserByIDSQL       = getUsersSQL + " WHERE id = $1"
	getUserByUsernameSQL = getUsersSQL + " WHERE username = $1"
	getUserByEmailSQL = getUsersSQL + " WHERE email = $1"
	getUsersCountSQL     = "SELECT COUNT(*) FROM users"
	updateUserSQL        = "UPDATE users SET username = $1, email = $2, password_hash = $3, is_active = $4 WHERE id = $5"
	deleteUserSQL        = "DELETE FROM users WHERE id = $1"
)

type (
	UserDomain interface {
		CreateUser(ctx context.Context, operations db.SQLOperations, user *models.User) error
		GetUserByID(ctx context.Context, operations db.SQLOperations, id int64) (*models.User, error)
		GetUserByUsername(ctx context.Context, operations db.SQLOperations, username string) (*models.User, error)
		GetUserByEmail(ctx context.Context, operations db.SQLOperations, email string) (*models.User, error)
		GetUsersCount(ctx context.Context, operations db.SQLOperations, filter *models.Filter) (int, error)
		GetUsers(ctx context.Context, operations db.SQLOperations, filter *models.Filter) ([]*models.User, error)
		DeleteUser(ctx context.Context, operations db.SQLOperations, id int64) error
	}

	userDomain struct{}
)

func NewUserDomain() UserDomain {
	return &userDomain{}
}

func (s *userDomain) CreateUser(
	ctx context.Context,
	operations db.SQLOperations,
	user *models.User,
) error {
	user.Touch()

	if user.IsNew() {
		err := operations.QueryRowContext(
			ctx, createUserSQL,
			user.Username,
			user.Email,
			user.PasswordHash,
			user.IsActive,
		).Scan(&user.ID)
		if err != nil {
			return apperr.NewDatabaseError(err).LogErrorMessage("create user query error: %v", err)
		}
		return nil
	}
	_, err := operations.ExecContext(
		ctx,
		updateUserSQL,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.IsActive,
		user.ID,
	)
	if err != nil {
		return apperr.NewDatabaseError(err).LogErrorMessage("update user query error: %v", err)
	}
	return nil
}

func (s *userDomain) GetUserByID(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) (*models.User, error) {
	row := operations.QueryRowContext(
		ctx,
		getUserByIDSQL,
		id,
	)
	return s.scanRow(row)
}

func (s *userDomain) GetUserByUsername(
	ctx context.Context,
	operations db.SQLOperations,
	username string,
) (*models.User, error) {

	row := operations.QueryRowContext(
		ctx,
		getUserByUsernameSQL,
		username,
	)
	return s.scanRow(row)
}

func (s *userDomain) GetUserByEmail(
	ctx context.Context,
	operations db.SQLOperations,
	email string,
) (*models.User, error) {
	row := operations.QueryRowContext(
		ctx,
		getUserByEmailSQL,
		email,
	)
	return s.scanRow(row)
}

func (s *userDomain) GetUsersCount(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) (int, error) {

	query, args := s.buildQuery(getUsersCountSQL, filter.NoPagination())
	rows := operations.QueryRowContext(
		ctx,
		query,
		args...,
	)

	var count int
	err := rows.Scan(&count)
	if err != nil {
		return 0, apperr.NewDatabaseError(err).LogErrorMessage("get users count query error: %v", err)
	}
	return count, nil
}

func (s *userDomain) GetUsers(
	ctx context.Context,
	operations db.SQLOperations,
	filter *models.Filter,
) ([]*models.User, error) {
	query, args := s.buildQuery(getUsersSQL, filter.NoPagination())
	rows, err := operations.QueryContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		return []*models.User{}, apperr.NewDatabaseError(err).LogErrorMessage("get users query error: %v", err)
	}
	defer rows.Close()
	users := make([]*models.User, 0)
	for rows.Next() {
		user, err := s.scanRow(rows)
		if err != nil {
			return []*models.User{}, err
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return []*models.User{}, apperr.NewDatabaseError(rows.Err()).LogErrorMessage("get users query error: %v", err)
	}
	return users, nil
}

func (s *userDomain) DeleteUser(
	ctx context.Context,
	operations db.SQLOperations,
	id int64,
) error {
	_, err := operations.ExecContext(
		ctx,
		deleteUserSQL,
		id,
	)
	if err != nil {
		return apperr.NewDatabaseError(err).LogErrorMessage("delete user query error: %v", err)
	}
	return nil
}

func (s *userDomain) scanRow(
	row db.RowScanner,
) (*models.User, error) {
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, apperr.NewDatabaseError(err).LogErrorMessage("scan user row error: %v", err)
	}
	return &user, nil
}

func (s *userDomain) buildQuery(
	query string,
	filter *models.Filter,
) (string, []interface{}) {
	args := make([]interface{}, 0)
	conditions := make([]string, 0)
	counter := utils.NewPlaceholder()

	if filter.Term != "" {
		textCols := []string{"username", "email"}
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
