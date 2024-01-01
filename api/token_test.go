package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/aradwann/eenergy/db/mock"
	db "github.com/aradwann/eenergy/db/store"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createTestSession(t *testing.T, user db.User, server *Server) (db.Session, string) {
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	require.NoError(t, err, "failed to create token")

	session := db.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "unknown",
		ClientIp:     "unknown",
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}
	return session, refreshToken
}

func TestRenewAccessTokenAPI(t *testing.T) {
	user, _ := randomUser(t)

	testCases := []struct {
		name          string
		buildStubs    func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "RenewAccessToken_OK",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				session, refreshToken := createTestSession(t, user, server)

				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "RenewAccessToken_Unauthorized",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				// Create a session with a blocked flag
				session, refreshToken := createTestSession(t, user, server)
				session.IsBlocked = true

				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "RenewAccessToken_ExpiredSession",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				// Create a session with an expired ExpiresAt
				session, refreshToken := createTestSession(t, user, server)
				session.ExpiresAt = time.Now().Add(-time.Hour) // Set to a past time

				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "RenewAccessToken_SessionNotFound",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				// Create a session with an expired ExpiresAt
				session, refreshToken := createTestSession(t, user, server)

				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(db.Session{}, sql.ErrNoRows)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "RenewAccessToken_MismatchToken",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				// Create a session with an expired ExpiresAt
				session, refreshToken := createTestSession(t, user, server)
				session.RefreshToken = "mismatchttttoken"
				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "RenewAccessToken_IncorrectSessionUser",
			buildStubs: func(t *testing.T, store *mockdb.MockStore, user db.User, server *Server) string {
				// Create a session with an expired ExpiresAt
				session, refreshToken := createTestSession(t, user, server)
				session.Username = "mismatchusername"
				// Ensure GetSession is called with the expected arguments
				store.EXPECT().
					GetSession(gomock.Any(), gomock.Eq(session.ID)).
					Times(1).
					Return(session, nil)

				return refreshToken
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, store)
			refreshToken := tc.buildStubs(t, store, user, server)
			body := gin.H{
				"refresh_token": refreshToken,
			}
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(body)
			require.NoError(t, err)

			url := "/tokens/renew_access"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
