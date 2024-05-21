package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	var session Session
	row := q.callStoredFunction(ctx, "create_session",
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
		arg.CreatedAt,
	)

	err := scanSessionFromRow(row, &session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	var session Session
	row := q.callStoredFunction(ctx, "get_session", id)
	err := scanSessionFromRow(row, &session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func scanSessionFromRow(row *sql.Row, session *Session) error {

	err := row.Scan(
		&session.ID,
		&session.Username,
		&session.RefreshToken,
		&session.UserAgent,
		&session.ClientIp,
		&session.IsBlocked,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	// Check for errors after scanning
	if err := row.Err(); err != nil {
		// Handle row-related errors
		log.Fatal(err)
		return err
	}
	if err != nil {
		// Check for a specific error related to the scan
		if errors.Is(err, ErrRecordNotFound) {
			// fmt.Println("No rows were returned.")
			return err
		} else {
			// Handle other scan-related errors
			log.Fatal(err)
			return err
		}
	}

	return nil
}
