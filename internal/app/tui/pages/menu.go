package pages

import (
	"fmt"

	"github.com/rivo/tview"

	"github.com/usa4ev/ghostorange/internal/app/model"
)

func (c *Constructor) menu() *tview.List {
	menu := tview.NewList()

	c.Logger.Debugf("focus on menu page")

	for i := 0; i < model.KeyLimit; i++ {
		c.Logger.Debugf("request count for %v",
			model.GetItemTitle(i))

		n, err := c.Adapter.Count(i)
		if err != nil {
			c.ShowMessage(err.Error(), KeyMenu)
			c.Logger.Errorf("failed to count items: %v",
				err)
		}

		title := model.GetItemTitle(i)

		var pageKey string
		switch i {
		case model.KeyCredentials:
			pageKey = KeyCredentials
		case model.KeyText:
			pageKey = KeyText
		case model.KeyCards:
			pageKey = KeyCards
		case model.KeyBinary:
			pageKey = KeyBinary
		}

		c.Logger.Debugf("Adding menu item %v", i)

		menu.AddItem(fmt.Sprintf("%v (%v)", title, n),
			"",
			rune(49+i),
			func() {
				c.Build(pageKey)
				c.Pages.SwitchToPage(pageKey)
			})
	}

	return menu
}
