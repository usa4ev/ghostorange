package tui

import (
	"github.com/rivo/tview"

	"ghostorange/cmd/client/tui/pages"
	"ghostorange/internal/app/provider"
)

type (
	Application struct {
		tviewApp *tview.Application
		provider provider.Provider
		pages    *tview.Pages
	}
)

func New(provider provider.Provider) *Application {

	builder := pages.Constructor{
		Provider: provider,
		Pages:    tview.NewPages(),
	}

	// creating ui pages
	builder.Build(pages.KeyLoginForm)
	builder.Build(pages.KeyRegistrationForm)
	builder.Build(pages.KeyMenu)
	builder.Build(pages.KeyCredentials)

	app := tview.NewApplication().
		SetRoot(builder.GetPages().SwitchToPage(pages.KeyLoginForm), true).
		EnableMouse(true)

	return &Application{
		tviewApp: app,
		provider: provider,
		pages:    builder.GetPages(),
	}
}

func (app *Application) Run() error {
	return app.tviewApp.Run()
}
