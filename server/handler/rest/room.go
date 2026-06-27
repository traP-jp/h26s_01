package rest

import (
	"github.com/labstack/echo/v5"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traP-jp/h26s_01/server/api"
)

func (h *Handler) GetRooms(c *echo.Context) error {
	rooms, err := h.repo.ListRooms(c.Request().Context())
	if err != nil {
		return err
	}

	response := make([]api.Room, 0, len(rooms))
	for _, room := range rooms {
		response = append(response, api.Room{
			Id:      openapi_types.UUID(room.ID),
			Name:    room.Name,
			Members: []api.RoomMember{},
			Status:  api.RoomStatus(room.Status),
		})
	}

	return c.JSON(200, response)
}

func (h *Handler) PostRoom(c *echo.Context) error {
	return nil
}
