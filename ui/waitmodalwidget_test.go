// Copyright 2021 Jan Renken

// This file is part of MongoTUI.

// MongoTUI is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// MongoTUI is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with MongoTUI.  If not, see <http://www.gnu.org/licenses/>.

package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildSpinner_WithValidIndex_ReturnsSpinner(t *testing.T) {
	result := buildSpinner(0)

	assert.Equal(t, "[white]╭[white]─[blue]─[blue]╮\n[lightgrey]│  [lightblue]│\n[lightgrey]╰[darkgrey]─[darkgrey]─[lightblue]╯\n", result)
}