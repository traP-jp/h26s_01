package socketio

import (
	"encoding/json"

	"github.com/WillYingling/pubsub"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/kanjipool"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/zishang520/socket.io/servers/socket/v3"
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

	remainLives := kanjipool.MaxIncorrect - incorrect
	if incorrect >= kanjipool.MaxIncorrect {
		// TODO: game:end (cleared: false)
		h.handleGameEnd(s, roomID, round, remainLives, false)
	} else if correct >= kanjipool.MaxCorrect {
		// TODO: game:end (cleared: true)
		h.handleGameEnd(s, roomID, round, remainLives, true)
	} else {
		// TODO: round:started
	}

	return nil
}

func (h *Handler) handleGameEnd(s *socket.Socket, s_roomId socket.Room, round model.Round, remainLives int, isClear bool) error {
	h.repo.ChangeGameStatus(s.Request().Context(), round.GameID, "completed")
	roundwithResult, err := h.repo.GetAllRounds(s.Request().Context(), round.GameID)
	if err != nil {
		return err
	}

	totalTime, err := h.repo.CalcTotalTimeMs(s.Request().Context(), round.GameID)
	if err != nil {
		return err
	}

	apiRounds := make([]api.Round, len(roundwithResult))
	for i, r := range roundwithResult {
		strokes := make([]api.Stroke, len(r.Strokes))
		for j, st := range r.Strokes {
			strokes[j] = api.Stroke{
				DrawerId: st.DrawerID,
				X1:       st.X1,
				Y1:       st.Y1,
				X2:       st.X2,
				Y2:       st.Y2,
			}
		}
		apiRounds[i] = api.Round{
			ActualAnswer:  r.ActualAnswer,
			GuesserAnswer: r.GuesserAnswer,
			GuesserId:     r.GuesserID,
			Id:            openapi_types.UUID(r.ID),
			Strokes:       strokes,
			TimeMs:        int(r.TimeMs),
		}
	}

	gameEndEvent := api.GameEndEvent{
		Cleared:        isClear,
		RemainingLives: remainLives,
		TotalTimeMs:    totalTime,
		Rounds:         apiRounds,
	}

	s.To(s_roomId).Emit("game:end", gameEndEvent)
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
