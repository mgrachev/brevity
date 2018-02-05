package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	str := GenerateToken(6)
	assert.Equal(t, len(str), 6)
}
