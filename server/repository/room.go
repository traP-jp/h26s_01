package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/kanjipool"
	"github.com/traP-jp/h26s_01/server/model"
)

func (r *Repository) ListRooms(ctx context.Context) ([]model.Room, error) {
	var rooms []model.Room
	var members []model.RoomMember
	if err := r.db.SelectContext(ctx, &rooms, "SELECT id, name, status FROM rooms WHERE status = 'playing' OR status = 'waiting'"); err != nil {
		return nil, err
	}
	if err := r.db.SelectContext(ctx, &members, "SELECT room_id, user_id, is_ready, is_connected FROM room_members"); err != nil {
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

	// この段階ではまだjoinさせない

	room := model.Room{
		ID:      roomId,
		Name:    roomName,
		Status:  "waiting",
		Members: []model.RoomMember{},
	}

	return room, nil
}

func (r *Repository) JoinRoom(ctx context.Context, roomId uuid.UUID, userId string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO room_members (room_id, user_id) VALUES (?, ?)", roomId, userId)
	return err
}

func (r *Repository) SetUserReady(ctx context.Context, roomId uuid.UUID, userId string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE room_members SET is_ready = true WHERE room_id = ? AND user_id = ?", roomId, userId)
	return err
}

func (r *Repository) StartGame(ctx context.Context, roomId uuid.UUID, playerCount int) error {
	if _, err := r.db.ExecContext(ctx, "UPDATE rooms SET status = 'playing' WHERE id = ?", roomId); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(ctx, "INSERT INTO games (room_id) VALUES (?)", roomId); err != nil {
		return err
	}

	kanjies, err := kanjipool.SelectKanjies(playerCount)
	if err != nil {
		return err
	}
	for i, kanji := range kanjies {
		gameKanjiesId, err := uuid.NewV7()
		if err != nil {
			return err
		}
		if _, err = r.db.ExecContext(ctx, "INSERT INTO game_kanjies (id, game_id, `character`, kanji_order) VALUES (?, ?, ?, ?)", gameKanjiesId, roomId, kanji.Char, i+1); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) GetRoom(ctx context.Context, roomId uuid.UUID) (*model.Room, error) {
	var room model.Room
	var members []model.RoomMember

	if err := r.db.GetContext(ctx, &room, "SELECT id, name, status FROM rooms WHERE id = ?", roomId); err != nil {
		return nil, err
	}

	if err := r.db.SelectContext(ctx, &members, "SELECT room_id, user_id, is_ready, is_connected FROM room_members WHERE room_id = ?", roomId); err != nil {
		return nil, err
	}

	room.Members = members

	return &room, nil
}

func (r *Repository) GetRoomMembersOrderedByGuesserOrder(ctx context.Context, roomId uuid.UUID) ([]model.RoomMember, error) {
	var members []model.RoomMember

	if err := r.db.SelectContext(ctx, &members, "SELECT room_id, user_id, is_ready, is_connected, guesser_order FROM room_members WHERE room_id = ? ORDER BY guesser_order ASC", roomId); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *Repository) ChangeGameStatus(ctx context.Context, roomId uuid.UUID, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE rooms SET status = '?' WHERE id = ?", status, roomId)
	return err
}
