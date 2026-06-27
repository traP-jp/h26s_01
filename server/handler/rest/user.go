package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *Handler) GetMe(c *echo.Context) error {
	user, err := h.getLoggedInUser(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}
