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

func (r *Repository) CountRoundResult(ctx context.Context, gameId uuid.UUID) (correct int, incorrect int, err error) {
	query := `SELECT
	SUM(rr.guesser_answer = gk.character) AS correct,
	SUM(rr.guesser_answer != gk.character) AS incorrect
	FROM rounds r
	JOIN round_results rr ON rr.round_id = r.id
	JOIN game_kanjies gk ON gk.id = r.kanji_id
	WHERE r.game_id = ?`

	var result struct {
		Correct   int `db:"correct"`
		Incorrect int `db:"incorrect"`
	}
	
	if err := r.db.GetContext(ctx, &result, query, gameId); err != nil {
		return 0, 0, err
	}
	return result.Correct, result.Incorrect, nil
}