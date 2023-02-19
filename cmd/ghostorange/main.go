package main

import (
	"log"

	"github.com/usa4ev/ghostorange/internal/app/server"
	"github.com/usa4ev/ghostorange/internal/app/srvconfig"
	"github.com/usa4ev/ghostorange/internal/app/storage"
)

func main() {
	cfg := srvconfig.New()

	strg, err := storage.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(cfg, strg)

	log.Fatal(srv.Run())
}
