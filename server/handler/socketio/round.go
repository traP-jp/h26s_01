package socketio

import (
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
