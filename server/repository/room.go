package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) ListRooms(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	var members []model.RoomMember
	if err := r.db.SelectContext(ctx, &rooms, "SELECT id, name, status FROM rooms WHERE status = 'playing' OR status = 'waiting'"); err != nil {
		return nil, err
	}
	if err := r.db.SelectContext(ctx, &members, "SELECT room_id, user_id, is_ready FROM room_members"); err != nil {
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

func (r *Repository) CreateRoom(ctx context.Context, roomName string, userId string) (model.Room, error) {
	roomId, err := uuid.NewV7()
	if err != nil {
		return model.Room{}, err
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO rooms (id, name, status) VALUES (?, ?, ?)", roomId, roomName, "waiting")
	if err != nil {
		return model.Room{}, err
	}

	_, err = r.db.ExecContext(ctx, "INSERT INTO room_members (room_id, user_id) VALUES (?, ?)", roomId, userId)
	if err != nil {
		return model.Room{}, err
	}

	roomMember := model.RoomMember{
		RoomID: roomId,
		UserID: userId,
	}

	room := model.Room{
		ID:      roomId,
		Name:    roomName,
		Status:  "waiting",
		Members: []model.RoomMember{roomMember},
	}

	return room, nil
}
