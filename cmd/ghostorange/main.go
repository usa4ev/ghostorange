package main

import (
	"log"

	"ghostorange/internal/app/server"
	"ghostorange/internal/app/srvconfig"
	"ghostorange/internal/app/storage"
)

func main(){
	cfg := srvconfig.New()
	
	strg,err := storage.New(cfg)
	if err != nil{
		log.Fatal(err)
	}

	srv := server.New(cfg, strg)
	
	log.Fatal(srv.Run())
}