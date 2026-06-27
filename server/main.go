package main

import (
	"log"

	"github.com/traP-jp/h26s_01/server/config"
	"github.com/traP-jp/h26s_01/server/database"
	"github.com/traP-jp/h26s_01/server/handler"
	"github.com/traP-jp/h26s_01/server/repository"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("runtime error: %v", err)
	}
}

func run() error {
	cfg := config.Load()
	db, err := database.Setup(cfg.MySQLConfig())

	if err != nil {
		return err
	}

	server := handler.NewServer(cfg, repository.New(db))

	return server.Start()
}
