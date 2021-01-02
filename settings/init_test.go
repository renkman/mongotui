package settings

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckKeyring_ReturnsClientRelatedResult(t *testing.T) {
	expected := os.Getenv("AGENT_ID") == ""

	result := checkKeyring()

	assert.Equal(t, expected, result)
}
