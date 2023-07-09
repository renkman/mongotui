package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultListAdd_AddsDocument(t *testing.T) {
	source := map[string]interface{}{
		"manufacturer": "Commodore",
	}

	read := func() map[string]interface{} {
		return source
	}

	resultList := CreateResultList()
	resultList.Add(read)
	result := resultList.Get()
	values := result[0].Get()

	assert.Len(t, result, 1)
	assert.Equal(t, source, values)
}

func TestResultListGet_AddsDocumentReturnsClone(t *testing.T) {
	source := map[string]interface{}{
		"manufacturer": "Commodore",
	}

	read := func() map[string]interface{} {
		return source
	}

	resultList := CreateResultList()
	resultList.Add(read)
	result := resultList.Get()
	values := result[0].Get()

	source["system"] = "Amiga 500"

	assert.Len(t, result, 1)
	assert.NotEqual(t, source, values)
}

func TestResultListGet_AddsDocumentsWithComplexTypes(t *testing.T) {
	source := []map[string]interface{}{
		map[string]interface{}{
			"manufacturer": "Commodore",
			"system":       "Amiga 500",
			"year":         1987,
			"cpu":          "Motorola 68000",
			"chipset":      []string{"Agnus", "Denise", "Paula"},
		},
		map[string]interface{}{
			"manufacturer": "Commodore",
			"system":       "Amiga 4000",
			"year":         1992,
			"cpu":          "Motorola 68040",
			"chipset":      []string{"Alice", "Lisa", "Paula"},
		},
	}

	resultList := CreateResultList()

	for i := range source {
		read := func() map[string]interface{} {
			return source[i]
		}
		resultList.Add(read)
	}

	result := resultList.Get()

	assert.Len(t, result, 2)
	assert.Equal(t, source[0], result[0].Get())
	assert.Equal(t, source[1], result[1].Get())
}
