package socketio

import (
	"context"

	"github.com/WillYingling/pubsub"
	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

const maxStrokeCount = 9

func (h *Handler) handleDrawStroke(s *socket.Socket, event api.DrawStrokeEvent) error {
	var roomID socket.Room
	for _, room := range s.Rooms().Keys() {
		if room != socket.Room(s.Id()) {
			roomID = room
			break
		}
	}

	roomUUID, err := uuid.Parse(string(roomID))
	if err != nil {
		return err
	}

	turnId, err := h.repo.GetTurnIDbyRoomID(context.Background(), roomUUID)
	if err != nil {
		return err
	}

	stroke := model.Stroke{
		TurnID: turnId,
		X1:     event.X1,
		Y1:     event.Y1,
		X2:     event.X2,
		Y2:     event.Y2,
	}

	if err = h.repo.SaveStroke(context.Background(), stroke); err != nil {
		return err
	}

	s.To(roomID).Emit("draw:stroke", event)

	round, err := h.repo.GetCurrentRoundByRoomID(context.Background(), roomUUID)
	if err != nil {
		return err
	}

	turnCount, err := h.repo.GetTurnCountByRoundID(context.Background(), round.ID)
	if err != nil {
		return err
	}

	if turnCount >= maxStrokeCount {
		return nil
	}

	currentTurn, err := h.repo.GetTurnByID(context.Background(), turnId)
	if err != nil {
		return err
	}

	members, err := h.repo.GetRoomMembersOrderedByGuesserOrder(context.Background(), roomUUID)
	if err != nil {
		return err
	}

	currentDrawerIndex := -1
	for i, member := range members {
		if member.UserID == currentTurn.DrawerID {
			currentDrawerIndex = i
			break
		}
	}

	if currentDrawerIndex == -1 {
		return nil
	}

	var nextDrawerID string
	for i := currentDrawerIndex + 1; i < len(members); i++ {
		if members[i].IsConnected && members[i].UserID != round.GuesserID {
			nextDrawerID = members[i].UserID
			break
		}
	}

	if nextDrawerID == "" {
		for i := 0; i < currentDrawerIndex; i++ {
			if members[i].IsConnected && members[i].UserID != round.GuesserID {
				nextDrawerID = members[i].UserID
				break
			}
		}
	}

	if nextDrawerID == "" {
		nextDrawerID = currentTurn.DrawerID
	}

	nextTurn := model.Turn{
		RoundID:   round.ID,
		TurnIndex: turnCount + 1,
		DrawerID:  nextDrawerID,
	}
	if err := h.repo.CreateTurn(context.Background(), &nextTurn); err != nil {
		return err
	}
	if err := h.repo.UpdateRoundCurrentTurn(context.Background(), round.ID, nextTurn.ID); err != nil {
		return err
	}

	turnStartedEvent := api.TurnStartedEvent{
		DrawerId:  nextTurn.DrawerID,
		TurnIndex: int(nextTurn.TurnIndex),
	}
	pubsub.Publish(context.Background(), turnStartedEvent)

	return nil
}
