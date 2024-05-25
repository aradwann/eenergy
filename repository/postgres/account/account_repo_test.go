package account

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aradwann/eenergy/entities"
	"github.com/aradwann/eenergy/repository/postgres/common"
	"github.com/aradwann/eenergy/repository/postgres/user"
	"github.com/aradwann/eenergy/util"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) *entities.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := user.CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testUserRepo.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

var testAccRepo AccountRepository
var testUserRepo user.UserRepository

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../..", ".env")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testAccRepo = NewAccountRepository(testDB)
	testUserRepo = user.NewUserRepository(testDB)
	os.Exit(m.Run())
}

func createRandomAccount(t *testing.T) *Account {
	u := CreateRandomUser(t)

	arg := CreateAccountParams{
		Owner:   u.Username,
		Balance: util.RandomAmount(),
		Unit:    util.RandomUnit(),
	}

	account, err := testAccRepo.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
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
	account2, err := testAccRepo.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
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

	account2, err := testAccRepo.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Unit, account2.Unit)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testAccRepo.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testAccRepo.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, common.ErrRecordNotFound.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount *Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testAccRepo.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)
	amount := int64(10)
	arg := AddAccountBalanceParams{
		ID:     account1.ID,
		Amount: amount,
	}

	updatedAccount, err := testAccRepo.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, account1.ID, updatedAccount.ID)
	require.Equal(t, account1.Owner, updatedAccount.Owner)
	require.Equal(t, account1.Balance+amount, updatedAccount.Balance)
	require.Equal(t, account1.Unit, updatedAccount.Unit)
	require.WithinDuration(t, account1.CreatedAt, updatedAccount.CreatedAt, time.Second)
}
