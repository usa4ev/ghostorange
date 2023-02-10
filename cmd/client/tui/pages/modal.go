package pages

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowError generates a new modal window with given message
// and sets the only button to switch focus to a page
// named in pageKey value.
func (c *Constructor) ShowError(message string, pageKey string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				c.Pages.SwitchToPage(pageKey)
			}
		})

	c.Pages.AddPage(KeyError, modal, false, false)
	c.Pages.SwitchToPage(KeyError)
}

// ShowInput generates a new page with input field
func (c *Constructor) ShowInput(message string, result *string, pageKey string) {
	input := tview.NewInputField().
	SetPlaceholder("file path").
	SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter{
			c.Pages.SwitchToPage(KeyInput)
		}
	}).
	SetChangedFunc(func(text string) {
		*result = text
	}).
		SetTitle(message)
		

	c.Pages.AddPage(KeyInput, input, false, false)
	c.Pages.SwitchToPage(KeyInput)
}