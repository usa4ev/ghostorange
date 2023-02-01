package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Contact struct {
	firstName   string
	lastName    string
	email       string
	phoneNumber string
	state       string
	business    bool
}

var contacts []Contact

var contactsList = tview.NewList().ShowSecondaryText(false)

var states = []string{"AK", "AL", "AR", "AZ", "CA", "CO", "CT", "DC", "DE", "FL", "GA",
	"HI", "IA", "ID", "IL", "IN", "KS", "KY", "LA", "MA", "MD", "ME",
	"MI", "MN", "MO", "MS", "MT", "NC", "ND", "NE", "NH", "NJ", "NM",
	"NV", "NY", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX",
	"UT", "VA", "VT", "WA", "WI", "WV", "WY"}

var demoApp = tview.NewApplication()

var contactText = tview.NewTextView()

var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(a) to add a new contact \n(q) to quit")

var demoForm = tview.NewForm()

var demoPages = tview.NewPages()

var flex = tview.NewFlex()

func notmain() {

	contactsList.SetSelectedFunc(func(index int, name string, second_name string, shortcut rune) {
		setConcatText(&contacts[index])
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().
			AddItem(contactsList, 0, 1, true).
			AddItem(contactText, 0, 4, false), 0, 6, false).
		AddItem(text, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 113 {
			demoApp.Stop()
		} else if event.Rune() == 97 {
			demoForm.Clear(true)
			addContactForm()
			demoPages.SwitchToPage("Add Contact")
		}
		return event
	})

	demoPages.AddPage("Menu", flex, true, true)
	demoPages.AddPage("Add Contact", demoForm, true, false)

	if err := demoApp.SetRoot(demoPages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func addContactForm() {
	contact := Contact{}

	demoForm.AddInputField("First Name", "", 20, nil, func(firstName string) {
		contact.firstName = firstName
	})

	demoForm.AddInputField("Last Name", "", 20, nil, func(lastName string) {
		contact.lastName = lastName
	})

	demoForm.AddInputField("Email", "", 20, nil, func(email string) {
		contact.email = email
	})

	demoForm.AddInputField("Phone", "", 20, nil, func(phone string) {
		contact.phoneNumber = phone
	})

	// states is a slice of state abbreviations. Code is in the repo.
	demoForm.AddDropDown("State", states, 0, func(state string, index int) {
		contact.state = state
	})

	demoForm.AddCheckbox("Business", false, func(business bool) {
		contact.business = business
	})

	demoForm.AddButton("Save", func() {
		contacts = append(contacts, contact)
		addContactList()
		demoPages.SwitchToPage("Menu")
	})

}

func addContactList() {
	contactsList.Clear()
	for index, contact := range contacts {
		contactsList.AddItem(contact.firstName+" "+contact.lastName, "", rune(49+index), nil)
	}
}

func setConcatText(contact *Contact) {
	contactText.Clear()
	text := contact.firstName + " " + contact.lastName + "\n" + contact.email + "\n" + contact.phoneNumber
	contactText.SetText(text)
}
