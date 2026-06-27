package repository

import (
	"context"

	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) ListRooms(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	var members []model.RoomMember
	if err := r.db.SelectContext(ctx, &rooms, "SELECT id, name, status FROM rooms"); err != nil {
		return nil, err
	}
	if err := r.db.SelectContext(ctx, &members, "SELECT room_id, user_id FROM room_members"); err != nil {
		return nil, err
	}

	// Create a map to associate room IDs with their members
	memberMap := make(map[string][]model.RoomMember)
	for _, member := range members {
		memberMap[member.RoomID.String()] = append(memberMap[member.RoomID.String()], member)
	}

	// Assign members to their respective rooms
	for i, room := range rooms {
		if roomMembers, ok := memberMap[room.ID.String()]; ok {
			rooms[i].Members = roomMembers
		}
	}	

	return rooms, nil
}

