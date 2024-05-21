package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/aradwann/eenergy/repository/postgres/common"
)

type VerifyEmail struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type EmailRepository interface {
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*VerifyEmail, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (*VerifyEmail, error)
}

type emailRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ EmailRepository = (*emailRepository)(nil)

func NewEmailRepository(db *sql.DB) EmailRepository {
	return &emailRepository{db: db}
}

type CreateVerifyEmailParams struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (r *emailRepository) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (*VerifyEmail, error) {
	verifyEmail := &VerifyEmail{}
	row := common.CallStoredFunction(ctx, r.db, "create_verify_email",
		arg.Username,
		arg.Email,
		arg.SecretCode,
	)

	if err := scanVerifyEmailFromRow(row, verifyEmail); err != nil {
		return nil, err
	}

	return verifyEmail, nil
}

type UpdateVerifyEmailParams struct {
	ID         int64  `json:"ID"`
	SecretCode string `json:"secret_code"`
}

func (r *emailRepository) UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (*VerifyEmail, error) {
	verifyEmail := &VerifyEmail{}
	params := []interface{}{
		arg.ID,
		arg.SecretCode,
	}
	row := common.CallStoredFunction(ctx, r.db, "update_verify_email", params...)

	// Execute the stored procedure and scan the results into the variables
	if err := scanVerifyEmailFromRow(row, verifyEmail); err != nil {
		return nil, err
	}

	return verifyEmail, nil
}

func scanVerifyEmailFromRow(row *sql.Row, verifyEmail *VerifyEmail) error {
	err := row.Scan(
		&verifyEmail.ID,
		&verifyEmail.Username,
		&verifyEmail.Email,
		&verifyEmail.SecretCode,
		&verifyEmail.IsUsed,
		&verifyEmail.CreatedAt,
		&verifyEmail.ExpiredAt,
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
