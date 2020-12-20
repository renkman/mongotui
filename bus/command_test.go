package bus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandGetMessage(t *testing.T) {
	command := Command{"Foo", func() string {
		return "Bar"
	}}

	got := command.getMessage()
	// if got != "Bar" {
	// 	t.Errorf("Wrong message %s instead of Bar", got)
	// }

	assert.Equal(t, "Bar", got)
}
