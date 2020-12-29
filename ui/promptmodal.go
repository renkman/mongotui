package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type PromptModal struct {
	*tview.Form
	*Widget
}

func createPromptModal(prompt string, ok func(), cancel func()) (*tview.Flex, *tview.Form) {
	passwordForm := tview.NewForm().
		AddPasswordField("", "", 30, '*', nil).
		AddButton("Ok", ok).
		AddButton("Cancel", cancel)

	frame := tview.NewFrame(passwordForm).
		AddText(prompt, true, tview.AlignCenter, tcell.ColorYellow)

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(frame, 20, 1, false).
			AddItem(nil, 0, 1, false), 40, 1, false).
		AddItem(nil, 0, 1, false)

	return modal, passwordForm
}

func CreatePromptModal(app *tview.Application, pages *tview.Pages, prompt string, ok func()) *PromptModal {
	modal, form := createPromptModal(prompt,
		ok,
		func() {
			pages.RemovePage("prompt")
		})
	widget := createWidget(modal, "prompt", app, pages)
	promptModal := PromptModal{form, widget}

	return &promptModal
}

func (p *PromptModal) GetPassword() string {
	return p.GetFormItem(0).(*tview.InputField).GetText()
}
