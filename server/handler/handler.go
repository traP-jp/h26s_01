package handler

import (
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/config"
	"github.com/traP-jp/h26s_01/server/handler/rest"
	"github.com/traP-jp/h26s_01/server/handler/socketio"
	"github.com/traP-jp/h26s_01/server/repository"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

type Server struct {
	config          *config.Config
	restAPIHandler  *rest.Handler
	socketIOHandler *socketio.Handler
}

func NewServer(config *config.Config, repo *repository.Repository) *Server {
	ioServer := socket.NewServer(nil, nil)

	return &Server{
		config:          config,
		restAPIHandler:  rest.NewHandler(repo, ioServer),
		socketIOHandler: socketio.NewHandler(repo, ioServer),
	}
}

func (s *Server) Start() error {
	e := echo.New()
	handler := s.restAPIHandler

	if err := s.socketIOHandler.Start(); err != nil {
		return err
	}

	api.RegisterHandlers(e, handler)
	e.Any("/socket.io/*", echo.WrapHandler(s.socketIOHandler.ServeHandler))

	return e.Start(s.config.AppAddr)
}
