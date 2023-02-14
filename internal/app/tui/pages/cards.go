package pages

import (
	"fmt"
	"time"

	"github.com/rivo/tview"

	"github.com/usa4ev/ghostorange/internal/app/model"
	"github.com/usa4ev/ghostorange/internal/pkg/argon2hash"
)

func (c *Constructor) cardsList() listGenerator {
	rflex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	tName := tview.NewTextView()
	tCardNo := tview.NewTextView()
	tCardHolder := tview.NewTextView()
	tExpires := tview.NewTextView()
	tComment := tview.NewTextView()

	tName.SetTitle("Name: ")
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

	rflex.AddItem(tName, 1, 0, false).
		AddItem(tCardNo, 1, 0, false).
		AddItem(tview.NewFlex().
			AddItem(tCardHolder, 0, 1, false).
			AddItem(tExpires, 0, 1, false), 1, 0, false).
		AddItem(tComment, 1, 0, false).
		AddItem(tview.NewButton("Show more").
			SetSelectedFunc(func() {
				c.Build(KeyFormCVV)
				c.Pages.SwitchToPage(KeyFormCVV)
			}), 1, 0, false)

	return listGenerator{
		btns:         buttons,
		detail:       rflex,
		Constructor:  c,
		addItemFunc:  addItemF,
		selectedFunc: selectedF,
	}
}

func (c *Constructor) cardsForm() *tview.Form {
	form := tview.NewForm()
	var cvv string

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
			AddInputField("CVV", cvv, 25, nil, func(text string) {
				cvv = text
			}).
			AddTextArea("Comment", item.Comment, 25, 3, 0, func(text string) {
				item.Comment = text
			}).
			AddButton("Save", func() {
				// Hash CVV code
				hash, err := argon2hash.GenerateFromPassword(cvv,
					argon2hash.DefaultParams())
				if err != nil {
					c.ShowMessage(err.Error(), KeyFormCards)
					return
				}

				item.CVVHash = hash

				c.Logger.Debugf("attemtimg to store card item: %v", item)

				if item.ID == "" {

					// Add data
					if err := c.Adapter.AddData(model.KeyCards, item); err != nil {
						c.ShowMessage(err.Error(), KeyFormCards)
						return
					}
				} else {
					if err := c.Adapter.UpdateData(model.KeyCards, item); err != nil {
						c.ShowMessage(err.Error(), KeyFormCards)
						return
					}
					c.CurItem = item
				}
				form.Clear(true)
				c.Pages.SwitchToPage(KeyCards)
			}).
			AddButton("Cancel", func() {
				form.Clear(true)
				c.Build(KeyCards)
				c.Pages.SwitchToPage(KeyCards)
			})
	})

	return form
}

func (c *Constructor) cardsCVVForm() *tview.Form {
	var cvv string
	return tview.NewForm().AddInputField("CVV", cvv, 3, nil, func(text string) {
		cvv = text
	}).
		AddButton("Back", func() {
			c.Pages.SwitchToPage(KeyCards)
		}).
		AddButton("Done", func() {
			if len(cvv) != 3 {
				c.ShowMessage("CVV code must be 3 charachters long", KeyFormCVV)
				return
			}

			item := c.CurItem.(model.ItemCard)

			c.Logger.Debugf("call GetCard() with id=%v; cvv=%v", item.ID, cvv)

			var err error
			item, err = c.Adapter.GetCard(item.ID, cvv)
			if err != nil {
				c.ShowMessage(err.Error(), KeyFormCVV)
				return
			}

			c.CurItem = item

			c.Build(KeyFormCards)
			c.Pages.SwitchToPage(KeyFormCards)
		})
}
