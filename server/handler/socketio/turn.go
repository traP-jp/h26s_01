package socketio

import (
	"context"

	"github.com/WillYingling/pubsub"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleTurnStartedEvent(socket *socket.Socket) error {
	ctx := context.Background()
	eventCh, unsubscribe := pubsub.SubscribeTo[api.TurnStartedEvent](ctx)

	socket.On("disconnect", func(args ...any) {
		unsubscribe()
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		case event := <-eventCh:
			socket.Emit("turn:started", event)
		}
	}
}
