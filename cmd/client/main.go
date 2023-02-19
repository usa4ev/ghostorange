package main

import (
	"fmt"
	"log"

	"go.uber.org/zap"

	"github.com/usa4ev/ghostorange/internal/app/adapter"
	"github.com/usa4ev/ghostorange/internal/app/tui"
	"github.com/usa4ev/ghostorange/internal/app/tui/clconfig"
)

func main() {
	cfg := clconfig.New()

	lgcfg := zap.NewDevelopmentConfig()
	lgcfg.OutputPaths[0] = cfg.LogPath()

	logger, _ := lgcfg.Build(zap.AddCaller())
	defer logger.Sync()

	sugar := logger.Sugar()

	adapter, err := adapter.New(cfg, sugar)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create provider: %v", err.Error()))
	}

	app := tui.New(adapter, sugar)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
