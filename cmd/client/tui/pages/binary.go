package pages

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) binaryList() listGenerator {
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tFName := tview.NewTextView()
	tFSize := tview.NewTextView()
	tComment := tview.NewTextView()
	tComment.SetTitle("Comment: ")

	buttons := make(map[string]func())

	buttons["Add"] = func() {
		c.CurItem = model.ItemBinary{}
		c.Build(KeyFormLoadBinary)
		c.Pages.SwitchToPage(KeyFormLoadBinary)
	}

	buttons["Edit"] = func() {
		if val, ok := c.CurItem.(model.ItemBinary); ok && val.ID != "" {
			c.Pages.SwitchToPage(KeyFormLoadBinary)
		}
	}

	var data []model.ItemBinary

	addItemF := func(val any, list *tview.List) error {
		var ok bool
		data, ok = val.([]model.ItemBinary)
		if !ok {
			return fmt.Errorf("got unexpected data type; expected: %v",
				model.GetItemTitle(model.KeyBinary))
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		return nil
	}

	selectedF := func(index int, name string, second_name string, shortcut rune) {

		item := data[index]
		tFName.Clear().SetText(item.Name + "." + item.Extention)
		tFSize.Clear().SetText(strconv.Itoa(item.Size) + " byte")
		tComment.Clear().SetText(item.Comment)

		c.CurItem = item

		c.Logger.Debugf("CurItem set: %v", c.CurItem)
	}

	rflex.AddItem(tFName, 1, 0, false).
		AddItem(tFSize, 1, 0, false).
		AddItem(tComment, 1, 0, false).
		AddItem(tview.NewButton("Save to file").
			SetSelectedFunc(func() {
				c.Build(KeyFormSaveBinary)
				c.Pages.SwitchToPage(KeyFormSaveBinary)
			}), 1, 0, false).
		SetBlurFunc(c.forgetCurItem)

	return listGenerator{
		btns:         buttons,
		detail:       rflex,
		Constructor:  c,
		addItemFunc:  addItemF,
		selectedFunc: selectedF,
	}
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
			AddInputField("Comment", item.Comment, 25, nil, func(text string) {
				item.Comment = text
			}).
			AddButton("Load from file", func() {
				c.Build(KeyFormLoadBinary)
			}).
			AddButton("Save", func() {
				// Send item to the server
				var err error
				if item.ID == "" {
					err = c.Adapter.AddData(model.KeyBinary, item)
				} else {
					err = c.Adapter.UpdateData(model.KeyBinary, item)
				}

				if err != nil {
					c.ShowMessage(fmt.Sprintf("Failed to save data:\n%v", err.Error()),
						KeyFormBinary)
					return
				}

				c.Build(KeyBinary)
				c.ShowMessage("Success!",
					KeyBinary)
				defer form.Clear(true)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyBinary)
			})
	})

	return form
}

func (c *Constructor) binaryLoadForm() *tview.Form {
	var path string

	item, ok := c.CurItem.(model.ItemBinary)
	if !ok {
		item = model.ItemBinary{}
	}

	return tview.NewForm().
		AddInputField("Comment", path, 25, nil, func(text string) {
			item.Comment = text
		}).
		AddInputField("Path", path, 25, nil, func(text string) {
			path = text
		}).
		AddButton("Load", func() {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				c.ShowMessage(fmt.Sprintf("Failed to read file:\n%v", err.Error()),
					KeyFormLoadBinary)
				return
			}

			st, err := os.Stat(path)
			if err != nil {
				c.ShowMessage(fmt.Sprintf("Failed to access file:\n%v", err.Error()),
					KeyFormLoadBinary)
				return
			}

			item.Data = base64.StdEncoding.EncodeToString(data)
			item.Name = filepath.Base(path)
			item.Extention = filepath.Ext(path)
			item.Size = int(st.Size())

			c.CurItem = item

			c.Build(KeyFormBinary)
			c.Pages.SwitchToPage(KeyFormBinary)
		}).
		AddButton("Cancel", func() {
			path = ""
			c.Pages.SwitchToPage(KeyBinary)
		})

}

func (c *Constructor) binarySaveForm() *tview.Form {
	var path string

	return tview.NewForm().
		AddInputField("Path", path, 25, nil, func(text string) {
			path = text
		}).
		AddButton("Save", func() {
			item, ok := c.CurItem.(model.ItemBinary)
			if !ok {
				return
			}

			if fileExists(path) {
				// ToDo: show modal dialogue
			}

			data, err := base64.StdEncoding.DecodeString(item.Data)
			if err != nil {
				c.ShowMessage(fmt.Sprintf("Failed to decode data:\n%v", err.Error()),
					KeyFormLoadBinary)
					return
			}

			if err := ioutil.WriteFile(path, data, 0644); err != nil {
				c.ShowMessage(fmt.Sprintf("Failed to save file:\n%v", err.Error()),
					KeyFormLoadBinary)
					return
			}

			c.Build(KeyBinary)
			c.Pages.SwitchToPage(KeyBinary)
		}).
		AddButton("Cancel", func() {
			path = ""
			c.Pages.SwitchToPage(KeyBinary)
		})

}

func fileExists(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	return false
}

func (c *Constructor) saveFile(item model.ItemBinary) error {

	return nil

}

func (c *Constructor) loadFile() (model.ItemBinary, error) {

	return model.ItemBinary{}, nil

}
