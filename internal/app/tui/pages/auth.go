package pages

import (
	"github.com/rivo/tview"

	"github.com/usa4ev/ghostorange/internal/app/model"
	"github.com/usa4ev/ghostorange/internal/app/tui/appinfo"
)

func (c *Constructor) loginForm() *tview.Form {
	creds := model.Credentials{}

	tAppInfo := tview.NewTextView().
		SetText(appinfo.AppInfo()).SetSize(2, 50)

	tAppInfo.SetTextAlign(tview.AlignLeft)

	loginForm := tview.NewForm().
		AddFormItem(tAppInfo).
		AddInputField("username", "", 25, nil, func(text string) {
			creds.Login = text
		}).
		AddPasswordField("password", "", 25, '*', func(text string) {
			creds.Password = text
		}).
		AddButton("Login", func() {
			c.Logger.Debugf("login attempt, user %v",
				creds.Login)
			if err := c.Adapter.Login(creds); err == nil {
				c.Logger.Debugf("successfull login, user %v",
					creds.Login)
				c.Build(KeyMenu)
				c.Pages.SwitchToPage(KeyMenu)
			} else {
				c.ShowMessage(err.Error(), KeyLoginForm)
			}
		}).
		AddButton("Register", func() {
			c.Pages.SwitchToPage(KeyRegistrationForm)
		})

	return loginForm
}

func (c *Constructor) regForm() *tview.Form {
	creds := model.Credentials{}

	regForm := tview.NewForm().
		AddInputField("username", creds.Login, 25, nil, func(text string) {
			creds.Login = text
		}).
		AddPasswordField("password", creds.Password, 25, '*', func(text string) {
			creds.Password = text
		}).
		AddButton("Back", func() {
			c.Pages.SwitchToPage(KeyLoginForm)
		}).
		AddButton("Register", func() {
			if err := c.Adapter.Register(creds); err == nil {
				c.Build(KeyMenu)
				c.Pages.SwitchToPage(KeyMenu)
				creds = model.Credentials{}
			} else {
				c.ShowMessage(err.Error(), KeyRegistrationForm)
			}
		})

	return regForm
}
