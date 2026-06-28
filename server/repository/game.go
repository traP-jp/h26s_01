package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) GetKanji(ctx context.Context, kanjiID uuid.UUID) (string, error) {
	var character string
	if err := r.db.GetContext(ctx, &character, "SELECT character FROM game_kanjies WHERE id = ?", kanjiID); err != nil {
		return "", err
	}
	return character, nil
}
