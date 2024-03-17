package worker

import (
	"context"
	"fmt"
	"testing"

	mockdb "github.com/aradwann/eenergy/db/mock"
	db "github.com/aradwann/eenergy/db/store"
	mockmail "github.com/aradwann/eenergy/mail/mock"
	"github.com/aradwann/eenergy/util"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		Role:           util.UserRole,
	}
	return
}

func TestProcessTaskSendVerifyEmail(t *testing.T) {

	user, _ := randomUser(t)

	testCases := []struct {
		name       string
		task       *asynq.Task
		buildStubs func(store *mockdb.MockStore, mailer *mockmail.MockEmailSender)
	}{
		{
			name: "OK",
			task: asynq.NewTask(TaskSendVerifyEmail, []byte(fmt.Sprintf(`{"username":"%s"}`, user.Username))),
			buildStubs: func(store *mockdb.MockStore, mailer *mockmail.MockEmailSender) {

				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(user, nil)
				// Use a matcher for dynamic secret code
				store.EXPECT().
					CreateVerifyEmail(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, cve db.CreateVerifyEmail) (db.VerifyEmail, error) {
						// Optionally, add more assertions here if needed to check the contents of cve
						require.Equal(t, cve.Email, user.Email)
						require.Equal(t, cve.Username, user.Username)
						return db.VerifyEmail{ID: 1, SecretCode: cve.SecretCode}, nil
					}).
					Times(1)

				mailer.EXPECT().
					SendEmail(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			storeCtrl := gomock.NewController(t)
			mockStore := mockdb.NewMockStore(storeCtrl)

			mailCtrl := gomock.NewController(t)
			mockMailer := mockmail.NewMockEmailSender(mailCtrl)

			tc.buildStubs(mockStore, mockMailer)

			// Initialize the processor with mocked dependencies
			processor := newTestTaskProcessor(t, mockStore, mockMailer)

			// Create a context and task for testing
			ctx := context.Background()

			// Call the method under test
			err := processor.ProcessTaskSendVerifyEmail(ctx, tc.task)

			// Assert expectations
			require.NoError(t, err)
		})
	}
}
