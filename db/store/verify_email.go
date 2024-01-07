package db

import (
	"context"
	"database/sql"
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

// func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
// 	var user User
// 	row := q.callStoredFunction(ctx, "get_user", username)
// 	err := scanUserFromRow(row, &user)
// 	if err != nil {
// 		return user, err
// 	}
// 	return user, nil
// }

// type UpdateUserParams struct {
// 	HashedPassword    sql.NullString `json:"hashed_password"`
// 	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
// 	FullName          sql.NullString `json:"fullname"`
// 	Email             sql.NullString `json:"email"`
// 	Username          string         `json:"username"`
// }

// func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
// 	var user User
// 	params := []interface{}{
// 		arg.Username,
// 		arg.HashedPassword,
// 		arg.PasswordChangedAt,
// 		arg.FullName,
// 		arg.Email,
// 	}
// 	row := q.callStoredFunction(ctx, "update_user", params...)

// 	// Execute the stored procedure and scan the results into the variables
// 	err := scanUserFromRow(row, &user)
// 	if err != nil {
// 		return User{}, err
// 	}

// 	return user, nil
// }

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
		if err == sql.ErrNoRows {
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
