package main

import (
	"fmt"
	"log"

	"ghostorange/cmd/client/tui"
	"ghostorange/internal/app/config"
	"ghostorange/internal/app/provider"
)

func main() {

	cfg := config.New()

	provider, err := provider.New(cfg)
	if err != nil{
		log.Fatal(fmt.Errorf("failed to create provider: %v", err.Error()))
	}

	app := tui.New(provider)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
