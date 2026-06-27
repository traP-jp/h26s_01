package rest

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/h26s_01/server/api"
)

func (h *Handler) GetMe(c *echo.Context) error {
	user, err := h.getLoggedInUser(c)
	if err != nil {
		return err
	}
	res := &api.User{
		Id: user.ID,
	}
	return c.JSON(http.StatusOK, res)
}
