package repository

import (
	"context"

	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetCurrentRoundByRoomID(ctx context.Context, roomID string) (*model.Round, error) {
	var round model.Round
	query := `
		SELECT r.id, r.game_id, r.round_index, r.current_turn_id, r.guesser_id, r.kanji_id, r.started_at
		FROM rounds r
		JOIN games g ON r.id = g.current_round_id
		WHERE g.room_id = ?`
	if err := r.db.GetContext(ctx, &round, query, roomID); err != nil {
		return nil, err
	}
	return &round, nil
}
