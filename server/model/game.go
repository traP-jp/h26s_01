package model

import "github.com/google/uuid"

type Game struct {
	RoomID         uuid.UUID     `db:"room_id"`
	CurrentRoundID uuid.NullUUID `db:"current_round_id"`
	Kanjies        []string      `db:"-"`
}

type GameKanji struct {
	ID         uuid.UUID `db:"id"`
	GameID     uuid.UUID `db:"game_id"`
	Character  string    `db:"character"`
	KanjiOrder uint8     `db:"kanji_order"`
}
