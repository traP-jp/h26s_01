package model

import "github.com/google/uuid"

type RoundResult struct {
	RoundID       uuid.UUID `db:"round_id"`
	GuesserAnswer string    `db:"guesser_answer"`
	TimeMs        uint32    `db:"time_ms"`
}
