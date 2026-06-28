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

// game:endでevent送信するための構造体 (strokeは別で送る)
type RoundWithResult struct {
    ID            uuid.UUID `db:"id"`
    GuesserID     string    `db:"guesser_id"`
    GuesserAnswer string    `db:"guesser_answer"`
    ActualAnswer  string    `db:"character"`
    TimeMs        uint32    `db:"time_ms"`
}

