package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aradwann/eenergy/repository/postgres/common"
	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type SessionRepository interface {
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error)
}

type sessionRepository struct {
	db *sql.DB
}

// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
var _ SessionRepository = (*sessionRepository)(nil)

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

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

func (r *sessionRepository) CreateSession(ctx context.Context, arg CreateSessionParams) (*Session, error) {
	session := &Session{}
	row := common.CallStoredFunction(ctx, r.db, "create_session",
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
		arg.CreatedAt,
	)

	if err := scanSessionFromRow(row, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	session := &Session{}
	row := common.CallStoredFunction(ctx, r.db, "get_session", id)
	if err := scanSessionFromRow(row, session); err != nil {
		return nil, err
	}

	return session, nil
}

func scanSessionFromRow(row *sql.Row, session *Session) error {
	if row == nil {
		return fmt.Errorf("row is nil")
	}
	if session == nil {
		return fmt.Errorf("session is nil")
	}
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
		if errors.Is(err, common.ErrRecordNotFound) {
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
