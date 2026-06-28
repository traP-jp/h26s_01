package socketio

import (
	"context"
	"errors"
	"time"

	"github.com/WillYingling/pubsub"
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleAnswerSubmit(socket *socket.Socket, event api.AnswerSubmitEvent) error {
	rooms := socket.Rooms()
	myID := string(socket.Id())
	var err error
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

	currentRound, err := h.repo.GetRoundByRoomID(context.Background(), roomID.String())
	if err != nil {
		return err
	}
	currentRoundUUID := currentRound.UUID

	currentTime := time.Now()

	if err := h.repo.SubmitAnswer(context.Background(), currentRoundUUID, currentTime, event.Answer); err != nil {
		return err
	}

	actualAnswer, err := h.repo.GetActualAnswer(context.Background(), currentRoundUUID)
	if err != nil {
		return err
	}

	roundAnswerEvent := api.RoundAnswerEvent{
		ActualAnswer:  actualAnswer,
		GuesserAnswer: event.Answer,
	}

	pubsub.Publish(context.Background(), roundAnswerEvent)

	return nil
}
