package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) SaveStroke(ctx context.Context, stroke model.Stroke) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	query := "INSERT INTO strokes (id, turn_id, x1, y1, x2, y2) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, id, stroke.TurnID, stroke.X1, stroke.Y1, stroke.X2, stroke.Y2)
	if err != nil {
		return err
	}
	return nil
}
