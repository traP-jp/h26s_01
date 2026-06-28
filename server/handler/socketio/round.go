package socketio

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/WillYingling/pubsub"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/zishang520/socket.io/servers/socket/v3"

	"github.com/google/uuid"
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
		h.handleRoundStarted(s, roomUUID)
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

func (h *Handler) handleRoundStartedEvent(socket *socket.Socket) error {

	ctx := socket.Request().Context()
	eventCh, unsubscribe := pubsub.SubscribeTo[api.RoundStartedEvent](ctx)

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

			socket.Emit("round:started", b)
		}
	}
}

func (h *Handler) handleRoundStarted(s *socket.Socket, roomUUID uuid.UUID) error {
	var round model.Round

	round.GameID = roomUUID

	members, err := h.repo.GetRoomMembersOrderedByGuesserOrder(s.Request().Context(), roomUUID)
	if err != nil {
		return err
	}

	kanjiesID, err := h.repo.GetKanjiesOrderByOrder(s.Request().Context(), roomUUID)
	if err != nil {
		return err
	}

	currentRound, err := h.repo.GetCurrentRoundByRoomID(s.Request().Context(), roomUUID)
	if err == nil {

		currentGuesserID := currentRound.GuesserID
		currentRoundIndex := currentRound.RoundIndex
		round.RoundIndex = currentRoundIndex + 1

		var lastOrder uint8 = 0
		if currentRoundIndex > 1 {
			for _, member := range members {
				if member.UserID == currentGuesserID {
					lastOrder = member.GuesserOrder
					break
				}
			}
		}

		var nextGuesserID string
		for _, member := range members {
			if member.GuesserOrder > lastOrder && member.IsConnected {
				nextGuesserID = member.UserID
				break
			}
		}

		if nextGuesserID == "" && len(members) > 0 {
			for _, member := range members {
				if member.IsConnected {
					nextGuesserID = member.UserID
					break
				}
			}
		}
		round.GuesserID = nextGuesserID

		currentKanjiID := currentRound.KanjiID

		if currentRoundIndex > 1 {
			for i, kanji := range kanjiesID {
				if kanji == currentKanjiID {
					round.KanjiID = kanjiesID[i+1]
					break
				}
			}
		}
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		} else {
			round.RoundIndex = 1
			if len(members) > 0 {
				round.GuesserID = members[0].UserID
			}
			if len(kanjiesID) > 0 {
				round.KanjiID = kanjiesID[0]
			}
		}
	}

	round.StartedAt = time.Now()

	if err := h.repo.CreateRound(s.Request().Context(), &round); err != nil {
		return err
	}

	kanji, err := h.repo.GetKanji(s.Request().Context(), round.KanjiID)
	if err != nil {
		return err
	}
	KanjiChar := kanji

	roundStartedEvent := api.RoundStartedEvent{
		GuesserId:  round.GuesserID,
		Kanji:      &KanjiChar,
		RoundIndex: int(round.RoundIndex),
	}

	pubsub.Publish(s.Request().Context(), roundStartedEvent)

	return nil
}
