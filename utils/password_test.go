package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = CheckPassword(password, hashedPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(10)
	err = CheckPassword(wrongPassword, hashedPassword)
	require.Equal(t, err.Error(), bcrypt.ErrMismatchedHashAndPassword.Error())
}
