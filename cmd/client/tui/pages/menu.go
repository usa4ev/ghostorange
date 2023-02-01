package pages

import (
	"strconv"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) menu() *tview.List {
	menu := tview.NewList()

	menu.SetFocusFunc(func() {
		for i := 0; i < model.KeyLimit; i++ {
			n, err := c.Provider.Count(i)
			if err != nil {
				// ToDo: handle error
			}

			title := model.GetItemTitle(i)

			var f func()
			if i == model.KeyCredentials{
				f = func(){
					c.Pages.SwitchToPage(KeyCredentials)
				}
			}

			menu.AddItem(title+" ("+strconv.Itoa(n)+")", "", rune(49+i),
			f)
		}
	})

	// regForm := tview.NewForm().
	// 	AddInputField("username", "", 25, nil, nil).
	// 	AddInputField("password", "", 25, nil, nil).
	// 	AddInputField("email", "", 25, nil, nil).
	// 	AddButton("Back", func() {
	// 		c.Pages.SwitchToPage(KeyLoginForm)
	// 	}).
	// 	AddButton("Register", func() {
	// 		c.Pages.SwitchToPage(KeyLoginForm) // ToDo: replace target form
	// 	})

	return menu
}

// func SelectedFunc(dataType int)func(){
// 	switch dataType{
// 	case KeyCredentials:
// 		return "Credentials"
// 	case KeyText:
// 		return "Text daata"
// 	case KeyBinary:
// 		return "Binary data"
// 	case KeyCard:
// 		return "Card info"
// 	}

// 	return ""
// }
