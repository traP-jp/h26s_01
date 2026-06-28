package socketio

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleGameReady(s *socket.Socket) error {
	slog.Info("Handling game:ready event", "socketID", s.Id())

	user, err := h.getLoggedInUser(s)
	if err != nil {
		slog.Error("Failed to get logged in user", "error", err)
		return err
	}

	slog.Info("User is ready", "userID", user.ID)

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

	slog.Info("Setting user ready", "roomID", roomID, "userID", user.ID)

	if err := h.repo.SetUserReady(context.Background(), roomID, user.ID); err != nil {
		return err
	}

	room, err := h.repo.GetRoom(context.Background(), roomID)

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
	h.io.To(socket.Room(roomID.String())).Emit("room:updated", &roomUpdatedEvent)

	slog.Info("Checking if all users are ready", "roomID", roomID)

	if h.isAllUsersReady(context.Background(), roomID) {
		slog.Info("All users ready, starting game", "roomID", roomID)
		if err := h.repo.StartGame(context.Background(), roomID, len(room.Members)); err != nil {
			return err
		}

		latestRoom, err := h.repo.GetRoom(context.Background(), roomID)
		if err != nil {
			return err
		}

		apiRoom := latestRoom.AsAPIRoom()
		roomUpdatedEvent := api.RoomUpdatedEvent{
			EventType: api.RoomUpdated,
			Room:      *apiRoom,
		}

		var roomListUpdatedEvent api.RoomListUpdatedEvent

		if err := roomListUpdatedEvent.FromRoomUpdatedEvent(roomUpdatedEvent); err != nil {
			return err
		}

		h.io.Emit("room_list:updated", &roomListUpdatedEvent)
		h.io.To(socket.Room(roomID.String())).Emit("room:updated", &roomUpdatedEvent)

		slog.Info("Game started, beginning first round", "roomID", roomID)
		h.handleRoundStarted(s, roomID)
	} else {
		slog.Info("Not all users ready yet", "roomID", roomID)
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
