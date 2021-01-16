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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type WaitModalWidget struct {
	*tview.Flex
	*tview.TextView
	*Widget
}

const (
	waitModal = "waitmodal"
	speed     = 100 * time.Millisecond
)

var highlighted = []string{"blue", "blue", "lightblue", "lightblue", "darkgrey", "darkgrey", "lightgrey", "lightgrey", "white", "white"}
var rotateOrder = []int{1, 2, 3, 7, 11, 10, 9, 8, 4, 0}
var spinnerRunes = []rune{'╭', '─', '─', '╮', '│', ' ', ' ', '│', '╰', '─', '─', '╯'}

func CreateWaitModalWidget(app *tview.Application, pages *tview.Pages, message string, ctx context.Context, cancel func()) *WaitModalWidget {
	spinner, modal, cancelButton := createWaitModalWidget(message, cancel)
	widget := createWidget(modal, waitModal, app, pages)

	pages.AddPage(waitModal, modal, true, true)
	app.SetFocus(cancelButton)
	waitModalWidget := WaitModalWidget{modal, spinner, widget}
	go waitModalWidget.rotateSpinner(ctx)
	return &waitModalWidget
}

func createWaitModalWidget(message string, cancel func()) (*tview.TextView, *tview.Flex, *tview.Button) {
	spinner := buildSpinner(0)
	spinnerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(spinner)
	spinnerTextView.SetBackgroundColor(tcell.ColorBlue)

	messageTextView := tview.NewTextView().SetText(message).SetWrap(true).SetTextAlign(tview.AlignCenter)
	messageTextView.SetBackgroundColor(tcell.ColorBlue)

	cancelButton := tview.NewButton("Cancel").
		SetSelectedFunc(cancel).
		SetLabelColor(tview.Styles.PrimaryTextColor)
	cancelButton.SetBackgroundColor(tview.Styles.ContrastBackgroundColor)

	innerModal := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(messageTextView, 60, 1, false).
			AddItem(nil, 0, 1, false),
			5, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(spinnerTextView, 3, 1, false).
				AddItem(nil, 0, 1, false),
				4, 1, false).
			AddItem(nil, 0, 1, false),
			4, 1, false).
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(cancelButton, 10, 1, false).
			AddItem(nil, 0, 1, false),
			1, 1, false).
		AddItem(nil, 0, 1, false)
	innerModal.SetBackgroundColor(tcell.ColorBlue).SetBorder(true)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(innerModal, 15, 1, false).
			AddItem(nil, 0, 1, false),
			80, 1, false).
		AddItem(nil, 0, 1, false)
	return spinnerTextView, modal, cancelButton
}

func (w *WaitModalWidget) rotateSpinner(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.Pages.RemovePage(waitModal)
			w.App.Draw()
			return
		default:
			for i := range rotateOrder {
				time.Sleep(speed)
				spinner := buildSpinner(i)
				w.SetText(spinner)
				w.App.Draw()
			}
		}
	}
}

func buildSpinner(index int) string {
	var builder strings.Builder
	for i, r := range spinnerRunes {
		ln := ""
		if (i+1)%4 == 0 {
			ln = "\n"
		}
		currentPosition := getPosition(i)
		if currentPosition == -1 {
			fmt.Fprintf(&builder, "%c%s", r, ln)
			continue
		}
		highlightedIndex := len(highlighted) - 1 - index + currentPosition
		if highlightedIndex >= len(highlighted) {
			highlightedIndex = highlightedIndex - len(highlighted)
		}
		color := highlighted[highlightedIndex]
		fmt.Fprintf(&builder, "[%s]%c%s", color, r, ln)
	}
	return builder.String()
}

func getPosition(index int) int {
	for pos, i := range rotateOrder {
		if i == index {
			return pos
		}
	}
	return -1
}
