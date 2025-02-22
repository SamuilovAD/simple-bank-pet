package db

import (
	"context"
	"github.com/SamuilovAD/simple-bank-pet/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	savedRandomAccount := createRandomAccount(t)
	fetchedRandomAccount, err := testQueries.GetAccount(context.Background(), savedRandomAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedRandomAccount)
	require.Equal(t, savedRandomAccount.Owner, fetchedRandomAccount.Owner)
	require.Equal(t, savedRandomAccount.Currency, fetchedRandomAccount.Currency)
	require.Equal(t, savedRandomAccount.Balance, fetchedRandomAccount.Balance)
}

func TestUpdateAccount(t *testing.T) {
	savedRandomAccount := createRandomAccount(t)
	args := UpdateAccountParams{
		ID:      savedRandomAccount.ID,
		Balance: util.RandomMoney(),
	}
	updatedRandomAccount, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedRandomAccount)
	require.Equal(t, savedRandomAccount.Owner, updatedRandomAccount.Owner)
	require.Equal(t, savedRandomAccount.Currency, updatedRandomAccount.Currency)
	require.Equal(t, args.Balance, updatedRandomAccount.Balance)
}

func TestDeleteAccount(t *testing.T) {
	savedRandomAccount := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), savedRandomAccount.ID)
	require.NoError(t, err)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Balance, account.Balance)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
