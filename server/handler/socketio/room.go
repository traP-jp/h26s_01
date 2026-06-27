package socketio

import (
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
	
	return nil
}
