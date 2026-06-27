package model

import "github.com/google/uuid"

type Room struct {
	ID      uuid.UUID    `db:"id"`
	Name    string       `db:"name"`
	Status  string       `db:"status"`
	Members []RoomMember `db:"-"`
}

type RoomMember struct {
	RoomID uuid.UUID `db:"room_id"`
	UserID string    `db:"user_id"`
}
