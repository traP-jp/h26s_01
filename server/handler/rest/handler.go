package rest

import (
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/traP-jp/h26s_01/server/repository"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

type Handler struct {
	repo *repository.Repository
	io   *socket.Server
}

func NewHandler(repo *repository.Repository, io *socket.Server) *Handler {
	return &Handler{
		repo: repo,
		io:   io,
	}
}

func (h *Handler) getLoggedInUser(c *echo.Context) (*model.User, error) {
	userID := c.Request().Header.Get("X-Forwarded-User")

	return h.repo.GetOrCreateUser(c.Request().Context(), userID)
}
