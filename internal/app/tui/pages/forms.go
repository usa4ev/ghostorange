package pages

import (
	"github.com/rivo/tview"
	"go.uber.org/zap"

	"github.com/usa4ev/ghostorange/internal/app/adapter"
	"github.com/usa4ev/ghostorange/internal/app/model"
)

const (
	// page keys
	KeyLoginForm        = "login form"
	KeyRegistrationForm = "registration form"
	KeyMenu             = "menu"
	KeyError            = "error"
	KeyInput            = "input"
	KeyCredentials      = "credentials"
	KeyFormCredentials  = "credentials form"
	KeyText             = "text"
	KeyFormText         = "text form"
	KeyCards            = "cards"
	KeyFormCards        = "cards form"
	KeyFormCVV          = "cvv input form"
	KeyBinary           = "binary"
	KeyFormBinary       = "binary form"
	KeyFormLoadBinary   = "binary load form"
	KeyFormSaveBinary   = "binary save form"
)

type (
	// Constructor creates new pages and add them to Pages.
	// Provider is requred to use in event handlers.
	Constructor struct {
		Adapter adapter.Adapter
		Pages   *tview.Pages
		CurItem any
		Logger  *zap.SugaredLogger
	}

	// listGenerator is builder for data type specific list-pages.
	listGenerator struct {
		*Constructor
		btns         map[string]func()
		detail       tview.Primitive
		key          string
		addItemFunc  func(any, *tview.List) error
		selectedFunc func(index int,
			name string,
			second_name string,
			shortcut rune)
	}
)

// Build creates new page depending on key and adds it
// to pages as a new page.
func (c *Constructor) Build(key string) {
	c.Pages.AddPage(key, c.buildPrimitive(key), true, false)
}

func (c *Constructor) buildPrimitive(key string) tview.Primitive {
	switch key {
	case KeyLoginForm:
		return c.loginForm()
	case KeyRegistrationForm:
		return c.regForm()
	case KeyMenu:
		return c.menu()
	case KeyCards, KeyBinary, KeyText, KeyCredentials:
		return c.BuildList(key)
	case KeyFormCredentials:
		return c.credForm()
	case KeyFormText:
		return c.textForm()
	case KeyFormCards:
		return c.cardsForm()
	case KeyFormCVV:
		return c.cardsCVVForm()
	case KeyFormBinary:
		return c.binaryForm()
	case KeyFormLoadBinary:
		return c.binaryLoadForm()
	case KeyFormSaveBinary:
		return c.binarySaveForm()
	default:
		return nil
	}
}

// GetPages returns pages. Expected to be called when all
// requred pages have been created.
func (c *Constructor) GetPages() *tview.Pages {
	return c.Pages
}

func (c *Constructor) forgetCurItem() {
	c.CurItem = nil
	c.Logger.Debugf("Current item set: %v", c.CurItem)
}

// BuildList creates a data type specific list-generator
// and newly generated list-page.
func (c *Constructor) BuildList(key string) tview.Primitive {
	lg := c.newListGenerator(key)
	return lg.build()
}

func (c *Constructor) newListGenerator(key string) listGenerator {
	var lg listGenerator

	switch key {
	case KeyCards:
		lg = c.cardsList()
	case KeyCredentials:
		lg = c.credList()
	case KeyBinary:
		lg = c.binaryList()
	case KeyText:
		lg = c.textList()
	}

	lg.key = key

	return lg
}

func (lg listGenerator) build() tview.Primitive {
	flex := tview.NewFlex()
	list := tview.NewList()
	lflex := tview.NewFlex().
		SetDirection(tview.FlexRow)

	// Create buttons
	menu := tview.NewFlex().
		AddItem(tview.NewButton("Back").
			SetSelectedFunc(func() {
				lg.forgetCurItem()
				lg.Build(KeyMenu)
				lg.Pages.SwitchToPage(KeyMenu)
			}), 0, 1, false)

	for k, v := range lg.btns {
		menu.AddItem(tview.NewButton(k).
			SetSelectedFunc(v), 0, 1, false)
	}

	// Fill the list
	val, err := lg.Adapter.GetData(listDataType(lg.key))
	if err != nil {
		lg.ShowMessage(err.Error(), KeyMenu)
		lg.Logger.Errorf("failed to get data: %v", err)
		return nil
	}

	// Add list rows
	if err = lg.addItemFunc(val, list); err != nil {
		lg.ShowMessage(err.Error(), KeyMenu)
	}

	list.SetSelectedFunc(lg.selectedFunc)

	// Compose the page
	lflex.AddItem(list, 0, 1, true).
		AddItem(menu, 1, 0, false)

	flex.AddItem(lflex, 0, 1, false).
		AddItem(lg.detail, 0, 1, false)

	return flex
}

func listDataType(Key string) int {
	switch Key {
	case KeyCards:
		return model.KeyCards
	case KeyCredentials:
		return model.KeyCredentials
	case KeyText:
		return model.KeyText
	case KeyBinary:
		return model.KeyBinary
	}

	return model.KeyLimit
}
