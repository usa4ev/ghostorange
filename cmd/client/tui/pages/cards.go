package pages

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	"ghostorange/internal/app/model"
)

func (c *Constructor) cardsList() listGenerator {
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tCardNo := tview.NewTextView()
	tCardHolder := tview.NewTextView()
	tExpires := tview.NewTextView()
	tComment := tview.NewTextView()

	tCardNo.SetTitle("Card no: ")
	tCardHolder.SetTitle("Cardholder: ")
	tExpires.SetTitle("Expires: ")
	tComment.SetTitle("Comment: ")

	buttons := make(map[string]func())

	buttons["Add"] = func() {
		c.forgetCurItem()
		c.Pages.SwitchToPage(KeyFormCards)
	}

	buttons["Edit"] = func() {
		if val, ok := c.CurItem.(model.ItemCard); ok && val.ID != "" {
			c.Pages.SwitchToPage(KeyFormCards)
		}
	}

	var data []model.ItemCard

	addItemF := func(val any, list *tview.List) error {
		var ok bool
		data, ok = val.([]model.ItemCard)
		if !ok {
			return fmt.Errorf("got unexpected data type; expected: %v",
				model.GetItemTitle(model.KeyCards))
		}

		for _, item := range data {
			list.AddItem(item.Name, item.Comment, '0', nil)
		}

		return nil
	}

	selectedF := func(index int, name string, second_name string, shortcut rune) {

		item := data[index]
		tCardNo.Clear().SetText(item.Number)
		tCardHolder.Clear().
			SetText(fmt.Sprint(item.CardholderName,
				item.CardholderSurename))
		tExpires.Clear().
			SetText(item.Exp.Format("01/06"))
		tComment.Clear().SetText(item.Comment)

		c.CurItem = item

		c.Logger.Debugf("CurItem set: %v", c.CurItem)
	}

	rflex.AddItem(tCardNo, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(tCardHolder, 0, 1, false).
			AddItem(tExpires, 0, 1, false), 1, 0, false).
		AddItem(tComment, 1, 0, false).
		SetBlurFunc(c.forgetCurItem)


	return listGenerator{
		btns: buttons,
		detail: rflex,
		Constructor: c,
		addItemFunc: addItemF,
		selectedFunc: selectedF,
	}
}

func (c *Constructor) cardsForm() *tview.Form {
	form := tview.NewForm()

	form.SetFocusFunc(func() {
		item := model.ItemCard{}

		if val, ok := c.CurItem.(model.ItemCard); ok {
			item = val
		}
		c.Logger.Debugf("filling the form using item %v", item)

		form.AddTextView("ID", item.ID, 50, 1, false, false).
			AddInputField("Name", item.Name, 25, nil, func(text string) {
				item.Name = text
			}).
			AddInputField("Number", item.Number, 25, nil, func(text string) {
				item.Number = text
			}).
			AddInputField("Cardholder name", item.CardholderName, 25, nil, func(text string) {
				item.CardholderName = text
			}).
			AddInputField("Cardholder surename", item.CardholderSurename, 25, nil, func(text string) {
				item.CardholderSurename = text
			}).
			AddInputField("Expiration month", item.Exp.Format("01/06"), 25, nil, func(text string) {
				exp, err := time.Parse("01/06", text)
				if err != nil {
					return
				}

				item.Exp = exp
			}).
			AddTextArea("Comment", item.Comment, 25, 3, 0, func(text string) {
				item.Comment = text
			}).
			AddButton("Save", func() {
				//ToDo: error handling
				if item.ID == "" {
					c.Adapter.AddData(model.KeyCards, item)
				} else {
					c.Adapter.UpdateData(model.KeyCards, item)
					c.CurItem = item
				}
				form.Clear(true)
				c.Pages.SwitchToPage(KeyCards)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Pages.SwitchToPage(KeyCards)
			})
	})

	return form
}
