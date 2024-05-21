package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"
)

type User struct {
	Username          string    `json:"username"`
	HashedPassword    string    `json:"hashed_password"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	IsEmailVerified   bool      `json:"is_email_verified"`
	Role              string    `json:"role"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (*User, error)
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

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	var user User
	slog.Info("CreateUser", slog.String("username", arg.Username))
	row := q.callStoredFunction(ctx, "create_user",
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email)

	err := scanUserFromRow(row, &user)
	if err != nil {
		slog.Error("error scaning the created user", slog.String("error message", err.Error()))
		return user, err
	}
	return user, nil
}

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	var user User
	row := q.callStoredFunction(ctx, "get_user", username)
	err := scanUserFromRow(row, &user)
	if err != nil {
		return user, err
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

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	var user User
	params := []interface{}{
		arg.Username,
		arg.HashedPassword,
		arg.PasswordChangedAt,
		arg.FullName,
		arg.Email,
		arg.IsEmailVerified,
	}
	row := q.callStoredFunction(ctx, "update_user", params...)

	// Execute the stored procedure and scan the results into the variables
	err := scanUserFromRow(row, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func scanUserFromRow(row *sql.Row, user *User) error {
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
		if errors.Is(err, ErrRecordNotFound) {
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

// callStoredFunction executes a stored function and returns a single row result.
func (r *userRepository) callStoredFunction(ctx context.Context, functionName string, params ...interface{}) *sql.Row {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)
	slog.Info("PostgreSQL function called",
		slog.String("function name", functionName),
		// slog.Any("params", params), // TODO: filter out sensitive info
		slog.String("SQL statement", sqlStatement),
	)

	return r.db.QueryRowContext(ctx, sqlStatement, params...)
}

// callStoredFunctionRows executes a stored function and returns multiple rows result.
func (r *userRepository) callStoredFunctionRows(ctx context.Context, functionName string, params ...interface{}) (*sql.Rows, error) {
	placeholders := generateParamPlaceholders(len(params))

	// Construct the SQL statement to call the stored function.
	sqlStatement := fmt.Sprintf(`SELECT * FROM %s(%s)`, functionName, placeholders)

	return r.db.QueryContext(ctx, sqlStatement, params...)
}

// generateParamPlaceholders generates placeholders for SQL parameters.
func generateParamPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 1; i <= count; i++ {
		placeholders[i-1] = fmt.Sprintf("$%d", i)
	}
	return strings.Join(placeholders, ", ")
}
