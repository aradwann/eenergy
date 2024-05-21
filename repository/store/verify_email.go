package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
)

type CreateVerifyEmail struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmail) (VerifyEmail, error) {
	var verifyEmail VerifyEmail
	row := q.callStoredFunction(ctx, "create_verify_email",
		arg.Username,
		arg.Email,
		arg.SecretCode,
	)

	err := scanVerifyEmailFromRow(row, &verifyEmail)
	if err != nil {
		return verifyEmail, err
	}
	return verifyEmail, nil
}

type UpdatVerifyEmailParams struct {
	ID         int64  `json:"ID"`
	SecretCode string `json:"secret_code"`
}

func (q *Queries) UpdatVerifyEmail(ctx context.Context, arg UpdatVerifyEmailParams) (VerifyEmail, error) {
	var verifyEmail VerifyEmail
	params := []interface{}{
		arg.ID,
		arg.SecretCode,
	}
	row := q.callStoredFunction(ctx, "update_verify_email", params...)

	// Execute the stored procedure and scan the results into the variables
	err := scanVerifyEmailFromRow(row, &verifyEmail)
	if err != nil {
		return verifyEmail, err
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
