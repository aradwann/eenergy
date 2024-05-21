package db

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomSession(t *testing.T) Session {

	user := createRandomUser(t)
	x := uuid.New()
	arg := CreateSessionParams{
		ID:           x,
		Username:     user.Username,
		RefreshToken: "",
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    time.Time{},
		CreatedAt:    time.Time{},
	}
	session, err := testStore.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, session.Username)
	require.Equal(t, arg.RefreshToken, session.RefreshToken)
	require.Equal(t, arg.UserAgent, session.UserAgent)
	require.Equal(t, arg.ClientIp, session.ClientIp)
	require.Equal(t, arg.IsBlocked, session.IsBlocked)
	require.Equal(t, arg.ClientIp, session.ClientIp)

	require.True(t, session.ExpiresAt.IsZero())
	require.NotZero(t, session.CreatedAt)
	return session
}

func TestCreateSession(t *testing.T) {
	createRandomSession(t)

}

func TestGetSession(t *testing.T) {
	session1 := createRandomSession(t)
	session2, err := testStore.GetSession(context.Background(), session1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, session2)

	require.Equal(t, session1.ID, session2.ID)
	require.Equal(t, session1.Username, session2.Username)
	require.Equal(t, session1.RefreshToken, session2.RefreshToken)
	require.Equal(t, session1.UserAgent, session2.UserAgent)
	require.Equal(t, session1.ClientIp, session2.ClientIp)
	require.Equal(t, session1.IsBlocked, session2.IsBlocked)

	require.WithinDuration(t, session1.CreatedAt, session2.CreatedAt, time.Second)
	require.WithinDuration(t, session1.ExpiresAt, session2.ExpiresAt, time.Second)
}
func TestGetSessionNotFound(t *testing.T) {
	session, err := testStore.GetSession(context.Background(), uuid.New())
	require.Equal(t, session, Session{})

	require.ErrorIs(t, err, ErrRecordNotFound)
}
