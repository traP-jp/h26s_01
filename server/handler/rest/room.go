package rest

import (
	"net/http"

	"github.com/WillYingling/pubsub"
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
		roomMembers := make([]api.RoomMember, len(room.Members))
		for i, member := range room.Members {
			roomMembers[i] = api.RoomMember{
				Id:      member.UserID,
				IsReady: member.IsReady,
			}
		}
		response = append(response, api.Room{
			Id:      openapi_types.UUID(room.ID),
			Name:    room.Name,
			Members: roomMembers,
			Status:  api.RoomStatus(room.Status),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateRoom(c *echo.Context) error {
	var body api.CreateRoomJSONRequestBody
	if err := c.Bind(&body); err != nil {
		return err
	}
	roomName := body.Name
	if roomName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "roomname is required")
	}

	user, err := h.getLoggedInUser(c)
	if err != nil {
		return err
	}

	room, err := h.repo.CreateRoom(c.Request().Context(), roomName, user.ID)
	if err != nil {
		return err
	}

	roomMembers := make([]api.RoomMember, len(room.Members))
	for i, member := range room.Members {
		roomMembers[i] = api.RoomMember{
			Id: member.UserID,
		}
	}
	response := api.Room{
		Id:      openapi_types.UUID(room.ID),
		Name:    room.Name,
		Members: roomMembers,
		Status:  api.RoomStatus(room.Status),
	}
	roomCreatedEvent := api.RoomCreatedEvent{
		EventType: api.RoomCreated,
		Room:      response,
	}

	var roomListUpdatedEvent api.RoomListUpdatedEvent

	if err := roomListUpdatedEvent.FromRoomCreatedEvent(roomCreatedEvent); err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	pubsub.Publish(c.Request().Context(), &roomListUpdatedEvent)

	return c.JSON(http.StatusCreated, response)
}
