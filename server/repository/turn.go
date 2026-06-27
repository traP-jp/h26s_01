package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) GetTurnIDbyRoomID(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error) {
	var turnId uuid.UUID
	query := "SELECT r.current_turn_id FROM games g JOIN rounds r ON r.id = g.current_round_id WHERE g.room_id = ?"	
	err := r.db.GetContext(ctx, &turnId, query, roomId)
	
	return turnId, err
}
