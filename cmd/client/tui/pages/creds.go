package pages

import (
	"fmt"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
	"ghostorange/internal/app/provider/httpp"
)

func (c *Constructor) credList() *tview.Flex {
	flex := tview.NewFlex()
	list := tview.NewList()
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tLogin := tview.NewTextView()
	tPassword := tview.NewTextView()
	tComment := tview.NewTextView()

	tLogin.SetTitle("Login: ")
	tPassword.SetTitle("Password: ")
	tComment.SetTitle("Comment: ")

	rflex.AddItem(tLogin, 1, 0, false).
		AddItem(tPassword, 1, 0, false).
		AddItem(tComment, 1, 0, false)

	flex.AddItem(list, 0, 1, true).
		AddItem(rflex, 0, 1, false)

	var data []model.ItemCredentials

	list.SetFocusFunc(func() {
		list.Clear()

		var err error

		val, err := c.Provider.GetData(model.KeyCredentials)
		if err != nil {
			// ToDo: handle error
			fmt.Printf("failed to get data: %v", err)
			return
		}

		var ok bool
		data,ok = val.([]model.ItemCredentials)
		if !ok{
			// ToDo: handle error
			_, ok = c.Provider.(httpp.Provider)
			fmt.Printf("%v\n", ok)
			fmt.Printf("data,ok = val.([]model.ItemCredentials); ok == %v", ok)
			return
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		list.SetSelectedFunc(
			func(index int, name string, second_name string, shortcut rune) {
				item := data[index]
				tLogin.Clear().SetText(item.Credentials.Login)
				tPassword.Clear().SetText(item.Credentials.Password)
				tComment.Clear().SetText(item.Comment)
			})

		
		
	})

	return flex
}
