package tui

import (
	"github.com/rivo/tview"
	"go.uber.org/zap"

	"github.com/usa4ev/ghostorange/internal/app/adapter"
	"github.com/usa4ev/ghostorange/internal/app/tui/pages"
)

type (
	Application struct {
		tviewApp *tview.Application
		adapter  adapter.Adapter
		pages    *tview.Pages
	}
)

func New(adapter adapter.Adapter, logger *zap.SugaredLogger) *Application {

	app := tview.NewApplication()

	builder := pages.Constructor{
		Adapter: adapter,
		Pages:   tview.NewPages(),
		Logger:  logger,
	}

	// Create ui pages
	builder.Build(pages.KeyLoginForm)
	builder.Build(pages.KeyRegistrationForm)
	builder.Build(pages.KeyMenu)

	builder.Build(pages.KeyFormCredentials)
	builder.Build(pages.KeyFormText)
	builder.Build(pages.KeyFormCards)
	builder.Build(pages.KeyFormBinary)

	app.SetRoot(builder.GetPages().
		SwitchToPage(pages.KeyLoginForm), true).
		EnableMouse(true)

	return &Application{
		tviewApp: app,
		adapter:  adapter,
		pages:    builder.GetPages(),
	}
}

func (app *Application) Run() error {
	return app.tviewApp.Run()
}
