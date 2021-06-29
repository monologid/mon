package mon_test

import (
	"github.com/monologid/mon"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJwt_EncryptDecryptShouldReturnSuccess(t *testing.T) {
	jwt := mon.NewJwt("secret", "HS256")

	data := map[string]interface{}{"name": "john doe"}
	token, err := jwt.Encrypt(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	tokenData, err := jwt.Decrypt(token)
	assert.NoError(t, err)
	assert.Equal(t, tokenData["name"], "john doe")
}

func TestNewJwt_DecryptShouldReturnError(t *testing.T) {
	jwt := mon.NewJwt("secret", "HS256")

	_, err := jwt.Decrypt("token")
	assert.Error(t, err)
}