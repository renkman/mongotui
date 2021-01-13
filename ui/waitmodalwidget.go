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
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type WaitModalWidget struct {
	*tview.Flex
	*tview.TextView
	*Widget
}

const (
	waitmodal string = "waitmodal"
	speed            = 200 * time.Millisecond
)

var highlighted = []string  {"blue", "blue", "blue", "blue", "darkgrey", "darkgrey", "lightgrey", "lightgrey", "white", "white"}
var rotateOrder = []int{1, 2, 3, 7, 11, 10, 9, 8, 4, 0}
var spinnerRunes = []rune{'╭', '─', '─', '╮', '│', ' ', ' ', '│', '╰', '─', '─', '╯'}

func CreateWaitModalWidget(app *tview.Application, pages *tview.Pages) *WaitModalWidget {
	spinner, modal := createWaitModalWidget()
	widget := createWidget(modal, waitmodal, app, pages)

	pages.AddPage(name, modal, true, true)
	app.SetFocus(modal)
	waitModalWidget := WaitModalWidget{modal, spinner, widget}
	go waitModalWidget.rotateSpinner()
	return &waitModalWidget
}

func createWaitModalWidget() (*tview.TextView, *tview.Flex) {
	spinner := buildSpinner(0)
	spinnerTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(spinner)

	spinnerTextView.SetBackgroundColor(tcell.ColorBlue)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(spinnerTextView, 5, 1, false).
			AddItem(nil, 0, 1, false), 6, 1, false).
		AddItem(nil, 0, 1, false)
	return spinnerTextView, modal
}

func (w *WaitModalWidget) rotateSpinner() {
	for i := 0; i < 10; i++ {
		for i, _ := range rotateOrder {
			time.Sleep(speed)
			spinner := buildSpinner(i)
			w.SetText(spinner)
			w.App.Draw()
		}
	}
}

func buildSpinner(index int) string {
	position := rotateOrder[index]
	var builder strings.Builder
	for i, r := range spinnerRunes {
		ln := ""
		if (i+1)%4 == 0 {
			ln = "\n"
		}
		offset := i - position
		color := highlighted[len(highlighted) - 1 - offset]
		fmt.Fprintf(&builder, "[%s]%c%s", color, r, ln)
	}
	return builder.String()
}
