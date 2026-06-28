package socketio

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/kanjipool"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleRoundEnd(s *socket.Socket) error {
	slog.Info("Handling round:end event", "socketID", s.Id())

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

	slog.Info("Getting current round", "roomID", roomUUID)

	round, err := h.repo.GetCurrentRoundByRoomID(context.Background(), roomUUID)
	if err != nil {
		return err
	}

	correct, incorrect, err := h.repo.CountRoundResult(context.Background(), round.GameID)
	if err != nil {
		return err
	}

	remainLives := kanjipool.MaxIncorrect - incorrect
	slog.Info("Round results", "correct", correct, "incorrect", incorrect, "remainLives", remainLives)

	if incorrect >= kanjipool.MaxIncorrect {
		// TODO: game:end (cleared: false)
		slog.Info("Game over - too many incorrect answers", "gameID", round.GameID)
		h.handleGameEnd(s, roomID, round, remainLives, false)
	} else if correct >= kanjipool.MaxCorrect {
		// TODO: game:end (cleared: true)
		slog.Info("Game cleared - all correct answers", "gameID", round.GameID)
		h.handleGameEnd(s, roomID, round, remainLives, true)
	} else {
		// TODO: round:started
		slog.Info("Starting next round", "gameID", round.GameID)
		h.handleRoundStarted(s, roomUUID)
	}

	return nil
}

func (h *Handler) handleGameEnd(s *socket.Socket, s_roomId socket.Room, round model.Round, remainLives int, isClear bool) error {
	h.repo.ChangeGameStatus(context.Background(), round.GameID, "completed")
	roundwithResult, err := h.repo.GetAllRounds(context.Background(), round.GameID)
	if err != nil {
		return err
	}

	totalTime, err := h.repo.CalcTotalTimeMs(context.Background(), round.GameID)
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

func (h *Handler) handleRoundStarted(s *socket.Socket, roomUUID uuid.UUID) error {
	var round model.Round

	round.GameID = roomUUID

	slog.Info("Starting round setup", "roomID", roomUUID)

	members, err := h.repo.GetRoomMembersOrderedByGuesserOrder(context.Background(), roomUUID)
	if err != nil {
		return err
	}
	slog.Info("Got room members", "count", len(members))

	kanjiesID, err := h.repo.GetKanjiesOrderByOrder(context.Background(), roomUUID)
	if err != nil {
		return err
	}
	slog.Info("Got kanjies", "count", len(kanjiesID))

	currentRound, err := h.repo.GetCurrentRoundByRoomID(context.Background(), roomUUID)
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

	slog.Info("Round configured", "roundIndex", round.RoundIndex, "guesserID", round.GuesserID, "kanjiID", round.KanjiID)

	round.StartedAt = time.Now()

	if err := h.repo.CreateRound(context.Background(), &round); err != nil {
		return err
	}
	slog.Info("Round created", "roundID", round.ID)

	if err := h.repo.UpdateGameCurrentRound(context.Background(), round.GameID, round.ID); err != nil {
		return err
	}
	slog.Info("Game current_round_id updated", "gameID", round.GameID, "roundID", round.ID)

	kanji, err := h.repo.GetKanji(context.Background(), round.KanjiID)
	if err != nil {
		return err
	}
	slog.Info("Kanji for round", "kanji", kanji)
	KanjiChar := kanji

	roundStartedEvent := api.RoundStartedEvent{
		GuesserId:  round.GuesserID,
		Kanji:      &KanjiChar,
		RoundIndex: int(round.RoundIndex),
	}
	slog.Info("Publishing round:started event", "event", roundStartedEvent)

	h.io.To(socket.Room(roomUUID.String())).Emit("round:started", roundStartedEvent)
	slog.Info("Published round:started event", "event", roundStartedEvent)

	var firstDrawerID string
	for _, member := range members {
		if member.UserID != round.GuesserID && member.IsConnected {
			firstDrawerID = member.UserID
			break
		}
	}

	if firstDrawerID != "" {
		turn := model.Turn{
			RoundID:   round.ID,
			TurnIndex: 1,
			DrawerID:  firstDrawerID,
		}
		if err := h.repo.CreateTurn(context.Background(), &turn); err != nil {
			return err
		}
		if err := h.repo.UpdateRoundCurrentTurn(context.Background(), round.ID, turn.ID); err != nil {
			return err
		}

		turnStartedEvent := api.TurnStartedEvent{
			DrawerId:  turn.DrawerID,
			TurnIndex: int(turn.TurnIndex),
		}
		h.io.To(socket.Room(roomUUID.String())).Emit("turn:started", turnStartedEvent)
		slog.Info("Published turn:started event", "event", turnStartedEvent)
	} else {
		slog.Warn("No first drawer found", "roundID", round.ID)
	}

	return nil
}
