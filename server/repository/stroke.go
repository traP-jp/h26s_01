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

func (r *Repository) GetAllStrokes (ctx context.Context, roundID uuid.UUID) ([]model.StrokeWithDrawerID, error) {
	var strokewithDrawerIDs []model.StrokeWithDrawerID
	query := 
	`SELECT t.drawer_id, s.x1, s.y1, s.x2, s.y2
	FROM strokes s
	JOIN turns t ON t.id = s.turn_id
	WHERE s.turn_id = ?`

	var turnIDs []uuid.UUID
	if err := r.db.SelectContext(ctx, &turnIDs, "SELECT id FROM turns WHERE round_id = ?", roundID); err != nil {
		return nil, err
	}

	var strokewithDrawerID model.StrokeWithDrawerID
	for _, turnID := range turnIDs {
		if err := r.db.SelectContext(ctx, &strokewithDrawerID, query, turnID); err != nil {
			return nil, err
		}
		strokewithDrawerIDs = append(strokewithDrawerIDs, strokewithDrawerID)
	}
	return strokewithDrawerIDs, nil
}
