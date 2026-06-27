package socketio

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleAnswerSubmit(socket *socket.Socket, event api.AnswerSubmitEvent) error {
	user, err := h.getLoggedInUser(socket)
	if err != nil {
		return err
	}

	rooms := socket.Rooms()
	myID := string(socket.Id())
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

	currentRoundID, err := h.repo.GetRoundByRoomID(socket.Request().Context(), roomID.String())
	if err != nil {
		return err
	}
	currentRoundUUID := currentRoundID.UUID

	current_time := time.Now()

	if err := h.repo.SubmitAnswer(socket.Request().Context(), currentRoundUUID, current_time, event.Answer); err != nil {
		return err
	}

	return nil
}
