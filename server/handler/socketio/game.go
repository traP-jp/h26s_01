package socketio

import (
	"context"
	"errors"

	"github.com/WillYingling/pubsub"
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleGameReady(s *socket.Socket) error {
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

	if err := h.repo.SetUserReady(s.Request().Context(), roomID, user.ID); err != nil {
		return err
	}

	room, err := h.repo.GetRoom(s.Request().Context(), roomID)

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
	pubsub.Publish(s.Request().Context(), roomUpdatedEvent)

	if h.isAllUsersReady(s.Request().Context(), roomID) {
		err := h.repo.StartGame(s.Request().Context(), roomID)
		if err != nil {
			return err
		}
		h.handleRoundStarted(s, roomID)
	}

	return nil
}

func (h *Handler) isAllUsersReady(ctx context.Context, roomID uuid.UUID) bool {
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
