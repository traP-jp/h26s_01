package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) SaveStroke(ctx context.Context, turnId uuid.UUID, x1, y1, x2, y2 float64) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	query := "INSERT INTO strokes (id, turn_id, x1, y1, x2, y2) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, id, turnId, x1, y1, x2, y2)
	if err != nil {
		return err
	}
	return nil
}
