package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetCurrentRoundByRoomID(ctx context.Context, roomID string) (*model.Round, error) {
	var round model.Round
	query := `
		SELECT r.id, r.game_id, r.round_index, r.current_turn_id, r.guesser_id, r.kanji_id, r.started_at
		FROM rounds r
		JOIN games g ON r.id = g.current_round_id
		WHERE g.room_id = ?`
	if err := r.db.GetContext(ctx, &round, query, roomID); err != nil {
		return nil, err
	}
	return &round, nil
}

func (r *Repository) CreateRound(ctx context.Context, gameID uuid.UUID, roundIndex uint8, guesserID string, kanjiID uuid.UUID) (uuid.UUID, error) {
	roundID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	query := "INSERT INTO rounds (id, game_id, round_index, guesser_id, kanji_id) VALUES (?, ?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, roundID, gameID, roundIndex, guesserID, kanjiID)
	if err != nil {
		return uuid.Nil, err
	}

	return roundID, nil
}

func (r *Repository) UpdateGameCurrentRound(ctx context.Context, gameID uuid.UUID, roundID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE games SET current_round_id = ? WHERE room_id = ?", roundID, gameID)
	return err
}

func (r *Repository) CreateTurn(ctx context.Context, roundID uuid.UUID, turnIndex uint8, drawerID string) (uuid.UUID, error) {
	turnID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}

	query := "INSERT INTO turns (id, round_id, turn_index, drawer_id) VALUES (?, ?, ?, ?)"
	_, err = r.db.ExecContext(ctx, query, turnID, roundID, turnIndex, drawerID)
	if err != nil {
		return uuid.Nil, err
	}

	return turnID, nil
}

func (r *Repository) UpdateRoundCurrentTurn(ctx context.Context, roundID uuid.UUID, turnID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE rounds SET current_turn_id = ? WHERE id = ?", turnID, roundID)
	return err
}
