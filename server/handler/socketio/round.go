package socketio

import (
	"encoding/json"

	"github.com/WillYingling/pubsub"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"

	"github.com/google/uuid"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

const (
	maxCorrect   = 9
	maxIncorrect = 3
)

func (h *Handler) handleRoundEnd(s *socket.Socket) error {
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

	round, err := h.repo.GetCurrentRoundByRoomID(s.Request().Context(), roomUUID)
	if err != nil {
		return err
	}

	correct, incorrect, err := h.repo.CountRoundResult(s.Request().Context(), round.GameID)
	if err != nil {
		return err
	}

	if incorrect >= maxIncorrect {
		// TODO: game:end (cleared: false)
	} else if correct >= maxCorrect {
		// TODO: game:end (cleared: true)
	} else {
		// TODO: round:started
	}

	return nil
}

func (h *Handler) handleRoundAnswer(socket *socket.Socket) error {
	ctx := socket.Request().Context()
	eventCh, unsubscribe := pubsub.SubscribeTo[api.RoundAnswerEvent](ctx)

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

			socket.Emit("round:answer", b)
		}
	}
}
