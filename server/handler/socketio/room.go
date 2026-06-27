package socketio

import (
	"github.com/WillYingling/pubsub"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleJoinRoom(s *socket.Socket, event api.RoomJoinEvent) error {
	user, err := h.getLoggedInUser(s)
	if err != nil {
		return err
	}
	err = h.repo.JoinRoom(s.Request().Context(), event.RoomId, user.ID)
	if err != nil {
		return err
	}

	s.Join(socket.Room(event.RoomId.String()))

	room, err := h.repo.GetRoom(s.Request().Context(), event.RoomId)

	if err != nil {
		return err
	}

	apiRoom := room.AsAPIRoom()
	roomUpdatedEvent := api.RoomUpdatedEvent{
		EventType: api.RoomUpdated,
		Room:      *apiRoom,
	}

	var roomListUpdatedEvent api.RoomListUpdatedEvent

	if err := roomListUpdatedEvent.FromRoomUpdatedEvent(roomUpdatedEvent); err != nil {
		return err
	}

	pubsub.Publish(s.Request().Context(), roomListUpdatedEvent)

	return nil
}
