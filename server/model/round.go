package model

import (
	"time"

	"github.com/google/uuid"
)

type Round struct {
	ID            uuid.UUID     `db:"id"`
	GameID        uuid.UUID     `db:"game_id"`
	RoundIndex    uint8         `db:"round_index"`
	CurrentTurnID uuid.NullUUID `db:"current_turn_id"`
	GuesserID     string        `db:"guesser_id"`
	KanjiID       uuid.UUID     `db:"kanji_id"`
	StartedAt     time.Time     `db:"started_at"`
}
