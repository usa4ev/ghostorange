package tui

import (
	"github.com/rivo/tview"
	"go.uber.org/zap"

	"ghostorange/cmd/client/tui/pages"
	"ghostorange/internal/app/adapter"
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

	// builder.Build(pages.KeyCredentials)
	builder.Build(pages.KeyFormCredentials)

	// builder.Build(pages.KeyText)
	builder.Build(pages.KeyFormText)

	// builder.Build(pages.KeyCards)
	builder.Build(pages.KeyFormCards)

	// builder.Build(pages.KeyBinary)
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
