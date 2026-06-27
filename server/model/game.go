package model

import "github.com/google/uuid"

type Game struct {
	RoomID         uuid.UUID     `db:"room_id"`
	CurrentRoundID uuid.NullUUID `db:"current_round_id"`
	Kanjies        []string      `db:"-"`
}
