package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/aradwann/eenergy/entities"
	"github.com/aradwann/eenergy/repository/postgres/common"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*entities.User, error)
	GetUser(ctx context.Context, username string) (*entities.User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (*entities.User, error)
}

type userRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ UserRepository = (*userRepository)(nil)

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"fullname"`
	Email          string `json:"email"`
}

func (r *userRepository) CreateUser(ctx context.Context, arg CreateUserParams) (*entities.User, error) {
	user := &entities.User{}
	slog.Info("CreateUser", slog.String("username", arg.Username))
	row := common.CallStoredFunction(ctx, r.db, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email)

	err := scanUserFromRow(row, user)
	if err != nil {
		slog.Error("error scaning the created user", slog.String("error message", err.Error()))
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUser(ctx context.Context, username string) (*entities.User, error) {
	user := &entities.User{}
	row := common.CallStoredFunction(ctx, r.db, "get_user", username)
	if err := scanUserFromRow(row, user); err != nil {
		return nil, err
	}
	return user, nil
}

type UpdateUserParams struct {
	HashedPassword    sql.NullString `json:"hashed_password"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	FullName          sql.NullString `json:"fullname"`
	Email             sql.NullString `json:"email"`
	Username          string         `json:"username"`
	IsEmailVerified   sql.NullBool   `json:"is_email_verified"`
}

func (r *userRepository) UpdateUser(ctx context.Context, arg UpdateUserParams) (*entities.User, error) {
	user := &entities.User{}
	params := []interface{}{
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
	}
	row := common.CallStoredFunction(ctx, r.db, "update_user", params...)

	// Execute the stored procedure and scan the results into the variables
	if err := scanUserFromRow(row, user); err != nil {
		return nil, err
	}
	return user, nil
}

func scanUserFromRow(row *sql.Row, user *entities.User) error {
	if row == nil {
		return fmt.Errorf("row is nil")
	}
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	err := row.Scan(
		&user.Username,
		&user.HashedPassword,
		&user.FullName,
		&user.Email,
		&user.PasswordChangedAt,
		&user.CreatedAt,
		&user.IsEmailVerified,
		&user.Role,
	)

	// Check for errors after scanning
	if err != nil {
		// Handle scan-related errors
		if errors.Is(err, common.ErrRecordNotFound) {
			// fmt.Println("No rows were returned.")
			return err
		} else {
			// Log and return other scan-related errors
			log.Printf("Error scanning row: %s", err)
			return err
		}
	}

	return nil
}
