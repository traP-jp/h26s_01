package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) GetOrCreateUser(ctx context.Context, userID string) (*model.User, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var user model.User

	if err = tx.
		QueryRowContext(ctx, "SELECT id FROM users WHERE id = ?", userID).
		Scan(&user.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = tx.
				QueryRowContext(ctx, "INSERT INTO users (id) VALUES (?) RETURNING id", userID).
				Scan(&user.ID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &user, tx.Commit()
}
