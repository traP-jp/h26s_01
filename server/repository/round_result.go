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

	var startTime time.Time
	err := r.db.GetContext(ctx, &startTime, "SELECT started_at FROM rounds WHERE id = ?", roundID)
	if err != nil {
		return err
	}

	timeMs := time.Since(startTime).Milliseconds()

	if _, err := r.db.ExecContext(ctx, "INSERT INTO round_results (round_id, guesser_answer, time_ms) VALUES (?, ?, ?)", roundResult.RoundID, roundResult.GuesserAnswer, timeMs); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetActualAnswer(ctx context.Context, roundID uuid.UUID) (string, error) {
	var actualAnswer string
	err := r.db.GetContext(ctx, &actualAnswer, "SELECT game_kanjies.character FROM rounds JOIN game_kanjies ON rounds.kanji_id = game_kanjies.id WHERE rounds.id = ?", roundID)
	if err != nil {
		return "", err
	}

	return actualAnswer, nil
}
