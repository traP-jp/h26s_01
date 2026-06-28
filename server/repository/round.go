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

func (r *Repository) CalcTotalTimeMs(ctx context.Context, roomID string) (int, error) {
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
