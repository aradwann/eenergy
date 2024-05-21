package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/aradwann/eenergy/util"

	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.NotZero(t, user.ID)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)

}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	_, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	updatedUser, err := testStore.GetUser(context.Background(), oldUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.FullName, newFullName)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Email, oldUser.Email)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	_, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	updatedUser, err := testStore.GetUser(context.Background(), oldUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.HashedPassword, oldUser.HashedPassword)
	require.Equal(t, updatedUser.Username, oldUser.Username)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(12)
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	_, err = testStore.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	updatedUser, err := testStore.GetUser(context.Background(), oldUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.Username, oldUser.Username)
	require.Equal(t, updatedUser.Email, oldUser.Email)
	require.Equal(t, updatedUser.Username, oldUser.Username)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := util.RandomString(12)
	newFullname := util.RandomOwner()
	newEmail := util.RandomEmail()

	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	_, err = testStore.UpdateUser(context.Background(), UpdateUserParams{
		ID: oldUser.ID,
		HashedPassword: sql.NullString{
			String: newHashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullname,
			Valid:  true,
		},
	})
	require.NoError(t, err)

	updatedUser, err := testStore.GetUser(context.Background(), oldUser.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updatedUser.HashedPassword, newHashedPassword)
	require.Equal(t, updatedUser.Email, newEmail)
	require.Equal(t, updatedUser.FullName, newFullname)

}
