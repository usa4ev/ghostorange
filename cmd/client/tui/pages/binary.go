package pages

import (
	"strconv"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) binaryList() *tview.Flex {
	flex := tview.NewFlex()
	list := tview.NewList()
	lflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tFName := tview.NewTextView()
	tFSize := tview.NewTextView()
	tComment := tview.NewTextView()

	tComment.SetTitle("Comment: ")

	lflex.AddItem(list, 0, 1, true).
		AddItem(tview.NewForm().
			AddButton("Add", func() {
				c.forgetCurItem()
				c.Pages.SwitchToPage(KeyFormBinary)
			}).
			AddButton("Edit", func() {
				if val, ok := c.CurItem.(model.ItemBinary); ok && val.ID != "" {
					c.Pages.SwitchToPage(KeyFormBinary)
				}
			}), 0, 1, false)

	rflex.AddItem(tFName, 1, 0, false).
		AddItem(tFSize, 1, 0, false).
		AddItem(tComment, 1, 0, false).
		AddItem(tview.NewButton("Save to file").
			SetSelectedFunc(func() {
				item, ok := c.CurItem.(model.ItemBinary)
				if !ok {
					return
				}
				if err := c.saveFile(item); err != nil {
					c.ShowError(err.Error(), KeyBinary)
				}
			}), 1, 0, false).
		SetBlurFunc(c.forgetCurItem)

	flex.AddItem(lflex, 0, 1, true).
		AddItem(rflex, 0, 1, false)

	var data []model.ItemBinary

	list.SetFocusFunc(func() {
		list.Clear()

		var err error

		val, err := c.Adapter.GetData(model.KeyBinary)
		if err != nil {
			// ToDo: handle error
			c.Logger.Errorf("failed to get data: %v", err)
			return
		}

		var ok bool
		data, ok = val.([]model.ItemBinary)
		if !ok {
			// ToDo: handle error
			return
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		c.Logger.Debugf("Binvary list set: %v", data)

		list.SetSelectedFunc(
			func(index int, name string, second_name string, shortcut rune) {

				item := data[index]
				tFName.Clear().SetText(item.Name + "." + item.Extention)
				tFSize.Clear().SetText(strconv.Itoa(item.Size) + " byte")
				tComment.Clear().SetText(item.Comment)

				c.CurItem = item

				c.Logger.Debugf("CurItem set: %v", c.CurItem)
			})

	})

	return flex
}

func (c *Constructor) binaryForm() *tview.Form {
	form := tview.NewForm()

	form.SetFocusFunc(func() {
		item := model.ItemBinary{}

		if val, ok := c.CurItem.(model.ItemBinary); ok {
			item = val
		}
		c.Logger.Debugf("filling the form using item %v", item)

		form.AddTextView("ID", item.ID, 50, 1, false, false).
			AddTextView("Size", strconv.Itoa(item.Size)+" byte", 50, 1, false, false).
			AddInputField("Name", item.Name, 25, nil, func(text string) {
				item.Name = text
			}).
			AddButton("Update from file", func() {
				//ToDo: error handling
				c.loadFile()
			}).
			AddButton("Save", func() {
				//ToDo: error handling
				if item.ID == "" {
					c.Adapter.AddData(model.KeyBinary, item)
				} else {
					c.Adapter.UpdateData(model.KeyBinary, item)
					c.CurItem = item
				}
				form.Clear(true)
				c.Pages.SwitchToPage(KeyBinary)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyBinary)
			})
	})

	return form
}

func (c *Constructor) saveFile(item model.ItemBinary) error {

	return nil

}

func (c *Constructor) loadFile() (model.ItemBinary, error) {

	return model.ItemBinary{}, nil

}
