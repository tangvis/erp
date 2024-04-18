package ecode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBusinessErrorCode(t *testing.T) {
	assert := require.New(t)

	assert.Equal(NewBusinessErrorCode(20, 187), 890120187)
}

func TestNewSystemErrorCode(t *testing.T) {
	assert := require.New(t)
	assert.Equal(NewSystemErrorCode(SystemReadWrite, 550), -8908550)
}
