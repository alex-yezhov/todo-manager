package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeTokenAndValidateToken(t *testing.T) {
	password := "12345"

	token, err := makeToken(password)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	assert.True(t, validateToken(token, password))
}

func TestValidateToken_WrongPassword(t *testing.T) {
	password := "12345"

	token, err := makeToken(password)
	require.NoError(t, err)

	assert.False(t, validateToken(token, "54321"))
}

func TestValidateToken_BadFormat(t *testing.T) {
	assert.False(t, validateToken("not-a-token", "12345"))
}

func TestPasswordHash(t *testing.T) {
	h1 := passwordHash("12345")
	h2 := passwordHash("12345")
	h3 := passwordHash("54321")

	require.NotEmpty(t, h1)

	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h1, h3)
}
