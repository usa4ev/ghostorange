package pages

import (
	//"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) loginForm() *tview.Form {
	creds := model.Credentials{}

	loginBtnAct := func() {
		if c.Provider.Login(creds) {
			// ToDo: successfull login logic
		} else {
			// ToDo: show error page
		}

		c.Pages.SwitchToPage(KeyMenu)
	}

	loginForm := tview.NewForm().
		AddInputField("username", "", 25, nil, func(text string) {
			creds.Login = text
		}).
		AddInputField("password", "", 25, nil, func(text string) {
			creds.Password = text
		}).
		AddButton("Login", loginBtnAct).
		AddButton("Register", func() {
			c.Pages.SwitchToPage(KeyRegistrationForm)
		})

	return loginForm
}

func (c *Constructor) regForm() *tview.Form {
	regForm := tview.NewForm().
		AddInputField("username", "", 25, nil, nil).
		AddInputField("password", "", 25, nil, nil).
		AddInputField("email", "", 25, nil, nil).
		AddButton("Back", func() {
			c.Pages.SwitchToPage(KeyLoginForm)
		}).
		AddButton("Register", func() {
			c.Pages.SwitchToPage(KeyLoginForm) // ToDo: replace target form
		})

	return regForm
}