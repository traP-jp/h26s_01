package model

import "github.com/google/uuid"

type Turn struct {
	ID        uuid.UUID `db:"id"`
	RoundID   uuid.UUID `db:"round_id"`
	TurnIndex uint8     `db:"turn_index"`
	DrawerID  string    `db:"drawer_id"`
}
