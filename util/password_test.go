package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)
	err = CheckPassword(password, hash1)
	require.NoError(t, err)
	wrongPassword := RandomString(5)
	err = CheckPassword(wrongPassword, hash1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)
}
