package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) SubmitAnswer(ctx context.Context, roundID uuid.UUID, currentTime time.Time, answer string) error {
	var roundResult model.RoundResult
	roundResult.RoundID = roundID
	roundResult.GuesserAnswer = answer
	roundResult.TimeMs = uint32(currentTime.Unix())

	var startTime uint32
	err := r.db.GetContext(ctx, &startTime, "SELECT started_at FROM rounds WHERE id = ?", roundID)
	if err != nil {
		return err
	}
	
	timeMs := roundResult.TimeMs - startTime

	if _, err := r.db.ExecContext(ctx, "INSERT INTO round_results (round_id, guesser_answer, time_ms) VALUES (?, ?, ?)", roundResult.RoundID, roundResult.GuesserAnswer, timeMs); err != nil {
		return err
	}

	return nil
}
