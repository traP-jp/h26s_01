package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetRoundResult(ctx context.Context, roundId uuid.UUID) (model.RoundResult, error) {
	var roundResult model.RoundResult
	if err := r.db.GetContext(ctx, &roundResult, "SELECT * FROM round_results WHERE round_id = ?", roundId); err != nil {
		return model.RoundResult{}, err
	}
	return roundResult, nil
}
