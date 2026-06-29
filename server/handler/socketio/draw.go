package socketio

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/traP-jp/h26s_01/server/model"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

const maxStrokeCount = 9

func (h *Handler) handleDrawStroke(s *socket.Socket, event api.DrawStrokeEvent) error {
	slog.Info("[draw:stroke] received", "socketID", s.Id(), "event", event)
	var roomID socket.Room
	for _, room := range s.Rooms().Keys() {
		if room != socket.Room(s.Id()) {
			roomID = room
			break
		}
	}
	slog.Info("[draw:stroke] roomID resolved", "roomID", roomID)

	roomUUID, err := uuid.Parse(string(roomID))
	if err != nil {
		slog.Error("[draw:stroke] failed to parse roomUUID", "error", err)
		return err
	}
	slog.Info("[draw:stroke] roomUUID parsed", "roomUUID", roomUUID)

	turnId, err := h.repo.GetTurnIDbyRoomID(context.Background(), roomUUID)
	if err != nil {
		slog.Error("[draw:stroke] failed to get turnID", "error", err)
		return err
	}
	slog.Info("[draw:stroke] turnID retrieved", "turnID", turnId)

	stroke := model.Stroke{
		TurnID: turnId,
		X1:     event.X1,
		Y1:     event.Y1,
		X2:     event.X2,
		Y2:     event.Y2,
	}
	slog.Info("[draw:stroke] stroke model created", "stroke", stroke)

	if err = h.repo.SaveStroke(context.Background(), stroke); err != nil {
		slog.Error("[draw:stroke] failed to save stroke", "error", err)
		return err
	}
	slog.Info("[draw:stroke] stroke saved")

	user, err := h.getLoggedInUser(s)
	if err != nil {
		slog.Error("[draw:stroke] failed to get logged in user", "error", err)
		return err
	}
	apiStroke := api.Stroke{
		DrawerId: user.ID,
		X1:       event.X1,
		Y1:       event.Y1,
		X2:       event.X2,
		Y2:       event.Y2,
	}
	s.To(roomID).Emit("draw:stroke", apiStroke)
	slog.Info("[draw:stroke] emitted to room", "roomID", roomID, "stroke", apiStroke)

	round, err := h.repo.GetCurrentRoundByRoomID(context.Background(), roomUUID)
	if err != nil {
		slog.Error("[draw:stroke] failed to get current round", "error", err)
		return err
	}
	slog.Info("[draw:stroke] current round retrieved", "roundID", round.ID, "roundIndex", round.RoundIndex)

	turnCount, err := h.repo.GetTurnCountByRoundID(context.Background(), round.ID)
	if err != nil {
		slog.Error("[draw:stroke] failed to get turn count", "error", err)
		return err
	}
	slog.Info("[draw:stroke] turn count retrieved", "turnCount", turnCount)

	if turnCount >= maxStrokeCount {
		slog.Info("[draw:stroke] max stroke count reached", "turnCount", turnCount, "maxStrokeCount", maxStrokeCount)
		return nil
	}

	currentTurn, err := h.repo.GetTurnByID(context.Background(), turnId)
	if err != nil {
		slog.Error("[draw:stroke] failed to get current turn", "error", err)
		return err
	}
	slog.Info("[draw:stroke] current turn retrieved", "currentTurn", currentTurn)

	members, err := h.repo.GetRoomMembersOrderedByGuesserOrder(context.Background(), roomUUID)
	if err != nil {
		slog.Error("[draw:stroke] failed to get room members", "error", err)
		return err
	}
	slog.Info("[draw:stroke] room members retrieved", "memberCount", len(members))

	currentDrawerIndex := -1
	for i, member := range members {
		if member.UserID == currentTurn.DrawerID {
			currentDrawerIndex = i
			break
		}
	}
	slog.Info("[draw:stroke] current drawer index found", "currentDrawerIndex", currentDrawerIndex, "currentDrawerID", currentTurn.DrawerID)

	if currentDrawerIndex == -1 {
		slog.Warn("[draw:stroke] current drawer not found in members")
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
		slog.Warn("[draw:stroke] no next drawer found, reusing current drawer", "nextDrawerID", nextDrawerID)
	} else {
		slog.Info("[draw:stroke] next drawer determined", "nextDrawerID", nextDrawerID)
	}

	nextTurn := model.Turn{
		RoundID:   round.ID,
		TurnIndex: turnCount + 1,
		DrawerID:  nextDrawerID,
	}
	if err := h.repo.CreateTurn(context.Background(), &nextTurn); err != nil {
		slog.Error("[draw:stroke] failed to create next turn", "error", err)
		return err
	}
	slog.Info("[draw:stroke] next turn created", "nextTurn", nextTurn)

	if err := h.repo.UpdateRoundCurrentTurn(context.Background(), round.ID, nextTurn.ID); err != nil {
		slog.Error("[draw:stroke] failed to update round current turn", "error", err)
		return err
	}
	slog.Info("[draw:stroke] round current turn updated", "roundID", round.ID, "nextTurnID", nextTurn.ID)

	turnStartedEvent := api.TurnStartedEvent{
		DrawerId:  nextTurn.DrawerID,
		TurnIndex: int(nextTurn.TurnIndex),
	}
	slog.Info("[draw:stroke] publishing turn:started event", "event", turnStartedEvent)
	h.io.To(roomID).Emit("turn:started", turnStartedEvent)
	slog.Info("[draw:stroke] turn:started event published")

	return nil
}
