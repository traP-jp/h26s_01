package model

import (
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
)

type Room struct {
	ID      uuid.UUID    `db:"id"`
	Name    string       `db:"name"`
	Status  string       `db:"status"`
	Members []RoomMember `db:"-"`
}

type RoomMember struct {
	RoomID       uuid.UUID `db:"room_id"`
	UserID       string    `db:"user_id"`
	IsReady      bool      `db:"is_ready"`
	GuesserOrder uint8     `db:"guesser_order"`
	IsConnected  bool      `db:"is_connected"`
}

func (r *Room) AsAPIRoom() *api.Room {
	roomMembers := make([]api.RoomMember, len(r.Members))

	for i, member := range r.Members {
		roomMembers[i] = api.RoomMember{
			Id:      member.UserID,
			IsReady: member.IsReady,
		}
	}

	return &api.Room{
		Id:      r.ID,
		Name:    r.Name,
		Members: roomMembers,
		Status:  api.RoomStatus(r.Status),
	}
}
