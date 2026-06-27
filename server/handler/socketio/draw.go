package socketio

import (
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleDrawStroke(s *socket.Socket, event api.DrawStrokeEvent) error {
	var roomID socket.Room
	for _, room := range s.Rooms().Keys() {
		if room != socket.Room(s.Id()) {
			roomID = room
			break
		}
	}
	
	roomUUID, err := uuid.Parse(string(roomID))
	if err != nil {
		return err
	}

	turnId, err := h.repo.GetTurnIDbyRoomID(s.Request().Context(), roomUUID)
	if err != nil {
		return err
	}

	if err = h.repo.SaveStroke(s.Request().Context(), turnId, event.X1, event.Y1, event.X2, event.Y2); err != nil {
		return err
	}

	s.To(roomID).Emit("draw:stroke", event)

	return nil
}
