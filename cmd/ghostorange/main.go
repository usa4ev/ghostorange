package main

import (
	"log"

	"ghostorange/internal/app/config"
	"ghostorange/internal/app/server"
)

func main(){
	cfg := config.New()
	srv := server.New(cfg)
	log.Fatal(srv.Run())
}