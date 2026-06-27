package main

import (
	"log"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/config"
	"github.com/traP-jp/h26s_01/server/database"
	"github.com/traP-jp/h26s_01/server/handler/rest"
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

	e := echo.New()
	handler := rest.NewHandler(repository.New(db))

	api.RegisterHandlers(e, handler)

	return e.Start(cfg.AppAddr)
}
