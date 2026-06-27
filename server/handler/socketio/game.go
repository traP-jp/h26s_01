package socketio

import (
	"context"

	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleGameReady(s *socket.Socket) error {
	user, err := h.getLoggedInUser(s)
	if err != nil {
		return err
	}

	rooms := s.Rooms()
	myID := string(s.Id())
	var roomID string

	if rooms != nil {
		for _, room := range rooms.Keys() {
			if string(room) != myID {
				roomID = string(room)
				break
			}
		}
	}
	if roomID == "" {
		return nil
	}

	if err := h.repo.SetUserReady(s.Request().Context(), roomID, user.ID); err != nil {
		return err
	}

	h.io.To(socket.Room(roomID)).Emit("game_ready", map[string]any{
		"user_id":  user.ID,
		"is_ready": true,
	})

	if h.isAllUsersReady(s.Request().Context(), roomID) {
		err := h.repo.StartGame(s.Request().Context(), roomID)
		if err != nil {
			return err
		}
		h.io.To(socket.Room(roomID)).Emit("game_start", map[string]any{
			"status":  "playing",
			"room_id": roomID,
		})
	}

	return nil
}

func (h *Handler) isAllUsersReady(ctx context.Context, roomID string) bool {
	room, err := h.repo.GetRoom(ctx, roomID)
	if err != nil {
		return false
	}
	if len(room.Members) < 2 {
		return false
	}
	for _, member := range room.Members {
		if !member.IsReady {
			return false
		}
	}
	return true
}
