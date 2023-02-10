package pages

import (
	"github.com/rivo/tview"

	"ghostorange/internal/app/provider"
)

const (
	// page keys
	KeyLoginForm        = "login form"
	KeyRegistrationForm = "registration form"
	KeyMenu = "menu"
	KeyCredentials = "credentials"
)

type (
	// Form struct {
	// 	*tview.Primitive
	// }

	// Constructor creates new pages and add them to Pages.
	// Provider is requred to use in event handlers.
	Constructor struct {
		Provider provider.Provider
		Pages    *tview.Pages
	}
)

// Build creates new page depending on key and adds it
// to pages as a new page.
func (c *Constructor) Build(key string){
	switch key {
	case KeyLoginForm:
		form := c.loginForm()
		c.Pages.AddPage(key, form, true, true)
	case KeyRegistrationForm:
		form := c.regForm()
		c.Pages.AddPage(key, form, true, true)
	case KeyMenu:
		list := c.menu()
		c.Pages.AddPage(key, list, true, true)
	case KeyCredentials:
		list := c.credList()
		c.Pages.AddPage(key, list, true, true)
	}
}

// GetPages returns pages. Expected to be called when all
// requred pages have been created.
func (c *Constructor) GetPages() *tview.Pages {
	return c.Pages
}

// // reset clears all input field on the form
// func (f *Form) reset() {
// 	r := f.GetFormItemCount()
// 	for i := 0; i < r; i++ {
// 		if inputField, ok := f.GetFormItem(i).(*tview.InputField); ok {
// 			inputField.SetText("")
// 		}
// 	}
// }
