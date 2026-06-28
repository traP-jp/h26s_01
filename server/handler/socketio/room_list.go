package socketio

import (
	"context"
	"encoding/json"

	"github.com/WillYingling/pubsub"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleRoomListUpdated(socket *socket.Socket) error {
	ctx := context.Background()
	eventCh, unsubscribe := pubsub.SubscribeTo[*api.RoomListUpdatedEvent](ctx)

	socket.On("disconnect", func(args ...any) {
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

			socket.Emit("room_list:updated", b)
		}
	}
}
