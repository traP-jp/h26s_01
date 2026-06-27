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
			if _, err = tx.ExecContext(ctx, "INSERT INTO users (id) VALUES (?) ON DUPLICATE KEY UPDATE id = id", userID); err != nil {
				return nil, err
			}
			user.ID = userID
		} else {
			return nil, err
		}
	}

	return &user, tx.Commit()
}
