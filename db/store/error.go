package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"
)

var ErrRecordNotFound = sql.ErrNoRows

var ErrUniqueViolation = &pgconn.PgError{
	Code: UniqueViolation,
}

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	// TODO: handle err conversion & logging
	fmt.Printf("errrrrr %#v", err)

	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
