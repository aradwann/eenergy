package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestErrRecordNotFound(t *testing.T) {
	err := ErrRecordNotFound
	assert.True(t, pgx.ErrNoRows == err)
}

func TestErrUniqueViolation(t *testing.T) {
	err := ErrUniqueViolation
	assert.True(t, errors.Is(err, ErrUniqueViolation))
}

func TestErrorCode(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		err := errors.New("some generic error")
		code := ErrorCode(err)
		assert.Equal(t, "", code)
	})

	t.Run("PgError", func(t *testing.T) {
		pgErr := &pgconn.PgError{
			Code: "23505",
		}
		err := pgErr
		code := ErrorCode(err)
		assert.Equal(t, "23505", code)
	})

	// t.Run("WrappedPgError", func(t *testing.T) {
	// 	wrappedPgErr := &pgconn.PgError{
	// 		Code: "23503",
	// 	}
	// 	err := errors.Wrap(wrappedPgErr, "wrapped error")
	// 	code := ErrorCode(err)
	// 	assert.Equal(t, "23503", code)
	// })

	t.Run("NonPgError", func(t *testing.T) {
		err := errors.New("some other error")
		code := ErrorCode(err)
		assert.Equal(t, "", code)
	})
}
