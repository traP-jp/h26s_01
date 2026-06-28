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

func (r *Repository) GetKanjiesOrderByOrder(ctx context.Context, gameID uuid.UUID) ([]uuid.UUID, error) {
	var kanjiesID []uuid.UUID
	if err := r.db.SelectContext(ctx, &kanjiesID, "SELECT id FROM game_kanjies WHERE game_id = ? ORDER BY kanji_order ASC", gameID); err != nil {
		return nil, err
	}
	return kanjiesID, nil
}
