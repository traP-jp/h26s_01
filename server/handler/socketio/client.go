package socketio

import (
	"errors"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleClientDisconnect(s *socket.Socket) error {
	user, err := h.getLoggedInUser(s)
	if err != nil {
		return err
	}

	rooms := s.Rooms()
	myID := string(s.Id())
	var roomID uuid.UUID

	if rooms != nil {
		for _, room := range rooms.Keys() {
			if string(room) != myID {
				roomID, err = uuid.Parse(string(room))
				if err != nil {
					return err
				}
				break
			}
		}
	}

	if roomID == uuid.Nil {
		return errors.New("roomID is empty")
	}

	clientDisconnectedEvent := api.ClientDisconnectedEvent{
		UserId: user.ID,
	}

	s.To(socket.Room(roomID.String())).Emit("client:disconnected", clientDisconnectedEvent)

	return nil
}
