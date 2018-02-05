package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestAppEnv(t *testing.T) {
	os.Setenv("APP_ENV", "foo")
	assert.Equal(t, AppEnv(), "foo")
	os.Setenv("APP_ENV", "")
}

func TestAppPort(t *testing.T) {
	os.Setenv("APP_PORT", "foo")
	assert.Equal(t, AppPort(), "foo")
	os.Setenv("APP_PORT", "")

	assert.Equal(t, AppPort(), defaultAppPort)
}

func TestAppDomain(t *testing.T) {
	os.Setenv("APP_DOMAIN", "foo")
	assert.Equal(t, AppDomain(), "foo")
	os.Setenv("APP_DOMAIN", "")

	assert.Equal(t, AppDomain(), defaultAppDomain)
}

func TestAppTokenLength(t *testing.T) {
	os.Setenv("APP_TOKEN_LENGTH", "7")
	assert.Equal(t, AppTokenLength(), 7)
	os.Setenv("APP_TOKEN_LENGTH", "")

	assert.Equal(t, AppTokenLength(), defaultAppTokenLength)
}
