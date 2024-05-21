package db

import (
	"context"
	"testing"
	"time"

	"github.com/aradwann/eenergy/util"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		OwnerUserID: user.ID,
		Balance:     util.RandomAmount(),
		Unit:        util.RandomUnit(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.OwnerUserID, account.OwnerUserID)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Unit, account.Unit)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerUserID, account2.OwnerUserID)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Unit, account2.Unit)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomAmount(),
	}

	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.OwnerUserID, account2.OwnerUserID)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Unit, account2.Unit)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		OwnerUserID: lastAccount.OwnerUserID,
		Limit:       5,
		Offset:      0,
	}

	accounts, err := testStore.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.OwnerUserID, account.OwnerUserID)
	}
}

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)
	amount := int64(10)
	arg := AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: amount,
	}

	updatedAccount, err := testStore.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, account1.ID, updatedAccount.ID)
	require.Equal(t, account1.OwnerUserID, updatedAccount.OwnerUserID)
	require.Equal(t, account1.Balance+amount, updatedAccount.Balance)
	require.Equal(t, account1.Unit, updatedAccount.Unit)
	require.WithinDuration(t, account1.CreatedAt, updatedAccount.CreatedAt, time.Second)
}
