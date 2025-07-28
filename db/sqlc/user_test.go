package db

import (
	"context"
	"github.com/SamuilovAD/simple-bank-pet/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testStore.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	args := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testStore.CreateUser(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.HashedPassword, user.HashedPassword)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)
	require.NotZero(t, user.CreatedAt)
	require.True(t, user.PasswordChangedAt.IsZero())

	return user
}

func TestUpdateUserFullName(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(),
		UpdateUserParams{
			Username: oldUser.Username,
			FullName: pgtype.Text{
				String: newFullName,
				Valid:  true,
			},
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, oldUser.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	newEmail := util.RandomOwner()
	newPassword := util.RandomOwner()
	newHashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	updatedUser, err := testStore.UpdateUser(context.Background(),
		UpdateUserParams{
			Username: oldUser.Username,
			FullName: pgtype.Text{
				String: newFullName,
				Valid:  true,
			},
			Email: pgtype.Text{
				String: newEmail,
				Valid:  true,
			},
			HashedPassword: pgtype.Text{
				String: newHashedPassword,
				Valid:  true,
			},
		},
	)

	require.NoError(t, err)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
}
