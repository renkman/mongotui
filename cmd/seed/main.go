// main.go
package main

import (
	"context"
	"time"

	"github.com/renkman/mongotui/testdata"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*11)
	defer cancel()

	testdata.Seed(ctx, "mongodb://localhost", 10000)
}
