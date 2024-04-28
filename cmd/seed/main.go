// main.go
package main

import (
	"context"
	"flag"
	"time"

	"github.com/renkman/mongotui/testdata"
)

var documentNumber int

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*11)
	defer cancel()

	initCommandLineArgs()
	flag.Parse()

	testdata.Seed(ctx, "mongodb://localhost", documentNumber)
}

func initCommandLineArgs() {
	flag.IntVar(&documentNumber,
		"n",
		10000,
		"Number of documents to create.")
}
