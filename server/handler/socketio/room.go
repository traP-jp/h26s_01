package socketio

import (
	"context"
	"errors"

	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

const MaxNumPlayers = 10

func (h *Handler) handleJoinRoom(s *socket.Socket, event api.RoomJoinEvent) error {
	user, err := h.getLoggedInUser(s)
	if err != nil {
		return err
	}
	currentRoom, err := h.repo.GetRoom(context.Background(), event.RoomId)
	if err != nil {
		return err
	}
	if len(currentRoom.Members) == MaxNumPlayers {
		return errors.New("The player limit has been reached.")
	}
	err = h.repo.JoinRoom(context.Background(), event.RoomId, user.ID)
	if err != nil {
		return err
	}

	s.Join(socket.Room(event.RoomId.String()))
	room, err := h.repo.GetRoom(context.Background(), event.RoomId)
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

	h.io.Emit("room_list:updated", &roomListUpdatedEvent)
	h.io.To(socket.Room(event.RoomId.String())).Emit("room:updated", &roomUpdatedEvent)

	return nil
}
