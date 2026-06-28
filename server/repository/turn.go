package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetTurnCountByRoundID(ctx context.Context, roundID uuid.UUID) (uint8, error) {
	var count int
	if err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM turns WHERE round_id = ?", roundID); err != nil {
		return 0, err
	}
	return uint8(count), nil
}

func (r *Repository) CreateTurn(ctx context.Context, turn *model.Turn) error {
	var err error
	turn.ID, err = uuid.NewV7()
	if err != nil {
		return err
	}
	query := "INSERT INTO turns (id, round_id, turn_index, drawer_id) VALUES (?, ?, ?, ?)"
	if _, err := r.db.ExecContext(ctx, query, turn.ID, turn.RoundID, turn.TurnIndex, turn.DrawerID); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTurnByID(ctx context.Context, turnID uuid.UUID) (model.Turn, error) {
	var turn model.Turn
	if err := r.db.GetContext(ctx, &turn, "SELECT id, round_id, turn_index, drawer_id FROM turns WHERE id = ?", turnID); err != nil {
		return model.Turn{}, err
	}
	return turn, nil
}

func (r *Repository) UpdateRoundCurrentTurn(ctx context.Context, roundID uuid.UUID, turnID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE rounds SET current_turn_id = ? WHERE id = ?", turnID, roundID)
	return err
}

func (r *Repository) GetTurnIDbyRoomID(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error) {
	var turnId uuid.UUID
	query := "SELECT r.current_turn_id FROM games g JOIN rounds r ON r.id = g.current_round_id WHERE g.room_id = ?"
	err := r.db.GetContext(ctx, &turnId, query, roomId)

	return turnId, err
}
