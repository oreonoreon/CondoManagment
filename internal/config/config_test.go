package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSessionSecret(t *testing.T) {
	t.Setenv("SESSION_SECRET", "")
	assert.Panics(t, func() { loadSessionSecret(true) }, "production без ключа должен падать")

	t.Setenv("SESSION_SECRET", "short")
	assert.Panics(t, func() { loadSessionSecret(true) }, "production с коротким ключом должен падать")

	long := strings.Repeat("k", minSessionSecretLen)
	t.Setenv("SESSION_SECRET", long)
	secret, generated := loadSessionSecret(true)
	assert.Equal(t, long, secret)
	assert.False(t, generated)

	t.Setenv("SESSION_SECRET", "")
	secret, generated = loadSessionSecret(false)
	assert.True(t, generated)
	assert.Len(t, secret, minSessionSecretLen*2)
}
