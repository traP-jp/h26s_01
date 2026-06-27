package model

import "github.com/google/uuid"

type Stroke struct {
	ID     uuid.UUID `db:"id"`
	TurnID uuid.UUID `db:"turn_id"`
	X1     float64   `db:"x1"`
	Y1     float64   `db:"y1"`
	X2     float64   `db:"x2"`
	Y2     float64   `db:"y2"`
}
