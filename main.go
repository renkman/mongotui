package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/renkman/mongotui/mongo"
	"github.com/rivo/tview"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		mongo.Disconnect(ctx)
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	app := tview.NewApplication()
	pages := tview.NewPages()

	CreateMainSreen(ctx, app, pages)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		e := fmt.Sprint(err)
		fmt.Print(e)
		panic(err)
	}
}
