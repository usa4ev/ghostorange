package pages

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/usa4ev/ghostorange/internal/app/model"
)

// Returns new list-page generator with all it needs to build
// a new page.
func (c *Constructor) credList() listGenerator {
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tLogin := tview.NewTextView()
	tPassword := tview.NewTextView()
	tComment := tview.NewTextView()

	tLogin.SetTitle("Login: ")
	tPassword.SetTitle("Password: ")
	tComment.SetTitle("Comment: ")

	buttons := make(map[string]func())

	buttons["Add"] = func() {
		c.forgetCurItem()
		c.Pages.SwitchToPage(KeyFormCredentials)
	}

	buttons["Edit"] = func() {
		if val, ok := c.CurItem.(model.ItemCredentials); ok && val.ID != "" {
			c.Pages.SwitchToPage(KeyFormCredentials)
		}
	}

	var data []model.ItemCredentials

	addItemF := func(val any, list *tview.List) error {
		var ok bool
		data, ok = val.([]model.ItemCredentials)
		if !ok {
			return fmt.Errorf("got unexpected data type; expected: %v",
				model.GetItemTitle(model.KeyCredentials))
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		return nil
	}

	selectedF := func(index int, name string, second_name string, shortcut rune) {

		item := data[index]
		tLogin.Clear().SetText(item.Credentials.Login)
		tPassword.Clear().SetText(item.Credentials.Password)
		tComment.Clear().SetText(item.Comment)

		c.CurItem = item

		c.Logger.Debugf("CurItem set: %v", c.CurItem)
	}

	rflex.AddItem(tLogin, 1, 0, false).
		AddItem(tPassword, 1, 0, false).
		AddItem(tComment, 1, 0, false).
		SetBlurFunc(c.forgetCurItem)

	return listGenerator{
		btns:         buttons,
		detail:       rflex,
		Constructor:  c,
		addItemFunc:  addItemF,
		selectedFunc: selectedF,
	}
}

// Returns data type specific add/edit form 
func (c *Constructor) credForm() *tview.Form {
	form := tview.NewForm()

	form.SetFocusFunc(func() {
		item := model.ItemCredentials{}

		if val, ok := c.CurItem.(model.ItemCredentials); ok {
			item = val
		}
		c.Logger.Debugf("filling the form using item %v", item)

		form.AddTextView("ID", item.ID, 50, 1, false, false).
			AddInputField("Name", item.Name, 25, nil, func(text string) {
				item.Name = text
			}).
			AddInputField("Login", item.Credentials.Login, 25, nil, func(text string) {
				item.Credentials.Login = text
			}).
			AddPasswordField("Password", item.Credentials.Password, 25, '*', func(text string) {
				item.Credentials.Password = text
			}).
			AddTextArea("Comment", item.Comment, 25, 3, 0, func(text string) {
				item.Comment = text
			}).
			AddButton("Save", func() {
				if item.ID == "" {
					if err := c.Adapter.AddData(model.KeyCredentials, item); err != nil {
						c.ShowMessage(err.Error(), KeyFormCredentials)
						return
					}
				} else {
					if err := c.Adapter.UpdateData(model.KeyCredentials, item); err != nil {
						c.ShowMessage(err.Error(), KeyFormCredentials)
						return
					}
					c.CurItem = item
				}
				form.Clear(true)
				c.Build(KeyCredentials)
				c.Pages.SwitchToPage(KeyCredentials)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyCredentials)
			})
	})

	return form
}
