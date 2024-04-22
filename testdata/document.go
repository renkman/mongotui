package testdata

import (
	"time"
)

type Document struct {
	Name        string
	Description string
	Created     time.Time
	Stuff       SubDocument
}

type ListDocument struct {
	Name    string
	Created time.Time
	List    []SubDocument
}

type SubDocument struct {
	Number int
	Text   string
}
