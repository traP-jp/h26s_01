package socketio

import (
	"encoding/json"

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
	pubsub.Publish(s.Request().Context(), roomUpdatedEvent)

	return nil
}

func (h *Handler) roomUpdatedEventHandler(s *socket.Socket) error {
	ctx := s.Request().Context()
	eventCh, unsubscribe := pubsub.SubscribeTo[api.RoomUpdatedEvent](ctx)

	s.On("disconnect", func(args ...any) {
		unsubscribe()
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-eventCh:
			b, err := json.Marshal(event)

			if err != nil {
				return err
			}

			s.Emit("room:updated", b)
		}
	}
}
