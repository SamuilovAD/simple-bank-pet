package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := RandomString(6)

	// Test successful hashing
	hash1, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	// Test that hashing the same password twice produces different hashes (due to salt)
	hash2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)

	// Test empty password
	emptyHash, err := HashPassword("")
	require.NoError(t, err)
	require.NotEmpty(t, emptyHash)

	// Test long password (bcrypt has a maximum length of 72 bytes)
	longPassword := RandomString(72)
	longHash, err := HashPassword(longPassword)
	require.NoError(t, err)
	require.NotEmpty(t, longHash)
}

func TestCheckPassword(t *testing.T) {
	password := RandomString(6)

	// Hash the password first
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	// Test correct password
	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	// Test wrong password
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// Test empty password
	err = CheckPassword("", hashedPassword)
	require.Error(t, err)

	// Test with empty hashed password
	err = CheckPassword(password, "")
	require.Error(t, err)

	// Test with invalid hash format
	err = CheckPassword(password, "invalid-hash-format")
	require.Error(t, err)
}
