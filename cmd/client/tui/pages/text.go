package pages

import (
	"fmt"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) textList() listGenerator{
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	txtView := tview.NewTextView()

	buttons := make(map[string]func())

	buttons["Add"] = func() {
		c.forgetCurItem()
		c.Pages.SwitchToPage(KeyFormText)
	}

	buttons["Edit"] = func() {
		if val, ok := c.CurItem.(model.ItemText); ok && val.ID != "" {
			c.Pages.SwitchToPage(KeyFormText)
		}
	}

	var data []model.ItemText

	addItemF := func(val any, list *tview.List) error {
		var ok bool
		data, ok = val.([]model.ItemText)
		if !ok {
			return fmt.Errorf("got unexpected data type; expected: %v",
				model.GetItemTitle(model.KeyText))
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		return nil
	}

	selectedF := func(index int, name string, second_name string, shortcut rune) {

		item := data[index]
		txtView.Clear().SetText(item.Text).SetTitle(item.Name)

		c.CurItem = item

		c.Logger.Debugf("CurItem set: %v", c.CurItem)
	}

	rflex.AddItem(txtView, 0, 1, false).
		SetBlurFunc(c.forgetCurItem)

	return listGenerator{
		btns:         buttons,
		detail:       rflex,
		Constructor:  c,
		addItemFunc:  addItemF,
		selectedFunc: selectedF,
	}
}

func (c *Constructor) textForm() *tview.Form {
	form := tview.NewForm()

	form.SetFocusFunc(func() {
		item := model.ItemText{}

		if val, ok := c.CurItem.(model.ItemText); ok {
			item = val
		}
		c.Logger.Debugf("filling the form using item %v", item)

		form.AddTextView("ID", item.ID, 50, 1, false, false).
			AddInputField("Name", item.Name, 25, nil, func(text string) {
				item.Name = text
			}).
			AddTextArea("Text", item.Text, 50, 10, 0, func(text string) {
				item.Text = text
			}).
			AddTextArea("Comment", item.Comment, 25, 3, 0, func(text string) {
				item.Comment = text
			}).
			AddButton("Save", func() {
				var err error
				if item.ID == "" {
					err = c.Adapter.AddData(model.KeyText, item)
					c.CurItem = item
				} else {
					err = c.Adapter.UpdateData(model.KeyText, item)
					c.CurItem = item
				}

				if err!= nil{
					c.ShowMessage(fmt.Sprintf("Failed to save text: %v", err.Error()),
						 KeyFormText)
					return
				}

				form.Clear(true)
				c.Build(KeyText)
				c.Pages.SwitchToPage(KeyText)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyText)
			})
	})

	return form
}
