package pages

import (
	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) textList() *tview.Flex {
	flex := tview.NewFlex()
	list := tview.NewList()
	lflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	txtView := tview.NewTextView()

	lflex.AddItem(list, 0, 1, true).
		AddItem(tview.NewForm().
			AddButton("Add", func() {
				c.forgetCurItem()
				c.Pages.SwitchToPage(KeyFormText)
			}).
			AddButton("Edit", func() {
				c.Pages.SwitchToPage(KeyFormText)
			}), 0, 1, false)

	rflex.AddItem(txtView, 0, 1, false).
		SetBlurFunc(c.forgetCurItem)

	flex.AddItem(lflex, 0, 1, true).
		AddItem(rflex, 0, 1, false)

	var data []model.ItemText

	list.SetFocusFunc(func() {
		list.Clear()

		var err error

		val, err := c.Adapter.GetData(model.KeyText)
		if err != nil {
			// ToDo: handle error
			c.Logger.Errorf("failed to get data: %v", err)
			return
		}

		var ok bool
		data, ok = val.([]model.ItemText)
		if !ok {
			// ToDo: handle error
			return
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		c.Logger.Debugf("Cred list set: %v", data)

		list.SetSelectedFunc(
			func(index int, name string, second_name string, shortcut rune) {

				item := data[index]
				txtView.Clear().SetText(item.Text).SetTitle(item.Name)

				c.CurItem = item

				c.Logger.Debugf("CurItem set: %v", c.CurItem)
			})

	})

	return flex
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
				//ToDo: error handling
				if item.ID == "" {
					c.Adapter.AddData(model.KeyText, item)
				} else {
					c.Adapter.UpdateData(model.KeyText, item)
					c.CurItem = item
				}
				form.Clear(true)
				c.Pages.SwitchToPage(KeyText)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyText)
			})
	})

	return form
}
