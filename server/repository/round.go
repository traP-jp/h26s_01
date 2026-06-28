package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetCurrentRoundByRoomID(ctx context.Context, roomID uuid.UUID) (model.Round, error) {
	var round model.Round
	query := "SELECT r.* FROM games g JOIN rounds r ON r.id = g.current_round_id WHERE g.room_id = ?"

	if err := r.db.GetContext(ctx, &round, query, roomID); err != nil {
		return model.Round{}, err
	}

	return round, nil
}

func (r *Repository) GetRoundByRoomID(ctx context.Context, roomID string) (uuid.NullUUID, error) {
	var game model.Game
	if err := r.db.GetContext(ctx, &game, "SELECT current_round_id FROM games WHERE room_id = ?", roomID); err != nil {
		return uuid.NullUUID{}, err
	}
	currentRoundID := game.CurrentRoundID

	return currentRoundID, nil
}

func (r *Repository) CreateRound(ctx context.Context, round *model.Round) error {
	var err error
	round.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO rounds (id, game_id, round_index, guesser_id, kanji_id, started_at) VALUES (?, ?, ?, ?, ?, ?) "
	if err := r.db.GetContext(ctx, &round, query, round.ID, round.GameID, round.RoundIndex, round.GuesserID, round.KanjiID, round.StartedAt); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetCurrentRoundIndexByRoomID(ctx context.Context, gameID uuid.UUID) (uint8, error) {
	var roundIndex uint8

	if err := r.db.GetContext(ctx, &roundIndex, "SELECT round_index FROM rounds WHERE game_id = ? ORDER BY round_index DESC LIMIT 1", gameID); err != nil {
		return 0, err
	}

	return roundIndex, nil
}
func (r *Repository) CalcTotalTimeMs(ctx context.Context, roomID uuid.UUID) (int, error) {
	var totalTime int
	var roundIDs []uuid.UUID
	if err := r.db.SelectContext(ctx, &roundIDs, "SELECT id FROM rounds WHERE game_id = ?", roomID); err != nil {
		return 0, err
	}

	for _, roundID := range roundIDs {
		roundResult, err := r.GetRoundResult(ctx, roundID)
		if err != nil {
			return 0, err
		}
		totalTime += int(roundResult.TimeMs)
	}
	return totalTime, nil
}

func (r *Repository) GetAllRounds(ctx context.Context, roomID uuid.UUID) ([]model.RoundWithResult, error) {
	var rounds []model.Round
	if err := r.db.SelectContext(ctx, &rounds, "SELECT * FROM rounds WHERE game_id = ?", roomID); err != nil {
		return nil, err
	}

	results := make([]model.RoundWithResult, 0, len(rounds))
	for _, round := range rounds {
		roundResult, err := r.GetRoundResult(ctx, round.ID) 
		if err != nil {
			return nil, err
		}

		kanji, err := r.GetKanji(ctx, round.KanjiID)
		if err != nil {
			return nil, err
		}

		strokes, err := r.GetAllStrokes(ctx, round.ID)
		if err != nil {
			return nil, err
		}
		
		results = append(results, model.RoundWithResult{
			ID: round.ID,
			GuesserID: round.GuesserID,
			GuesserAnswer: roundResult.GuesserAnswer,
			ActualAnswer: kanji,
			TimeMs: roundResult.TimeMs,
			Strokes: strokes,
		})
	}
	return results, nil
}
