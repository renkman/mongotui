package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckKeyring_ReturnsClientRelatedResult(t *testing.T) {
	result := checkKeyring()
	assert.True(t, result)
}
