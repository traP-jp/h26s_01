package rest

import (
	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/traP-jp/h26s_01/server/repository"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) getLoggedInUser(c *echo.Context) (*model.User, error) {
	userID := c.Request().Header.Get("X-Forwarded-User")

	return h.repo.GetOrCreateUser(c.Request().Context(), userID)
}
