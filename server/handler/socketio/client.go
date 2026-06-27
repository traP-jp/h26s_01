package socketio

import (
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/traP-jp/h26s_01/server/api"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

func (h *Handler) handleClientDisconnect(s *socket.Socket) error {
	ctx := s.Request().Context()
	user, err := h.getLoggedInUser(s)
	if err != nil {
		return err
	}

	rooms := s.Rooms()
	myID := string(s.Id())
	var roomID uuid.UUID

	if rooms != nil {
		for _, room := range rooms.Keys() {
			if string(room) != myID {
				roomID, err = uuid.Parse(string(room))
				if err != nil {
					return err
				}
				break
			}
		}
	}

	if roomID == uuid.Nil {
		return errors.New("roomID is empty")
	}

	clientDisconnectedEvent := api.ClientDisconnectedEvent{
		UserId: user.ID,
	}

	eventBytes, err := json.Marshal(clientDisconnectedEvent)
	if err != nil {
		return err
	}

	s.To(socket.Room(roomID.String())).Emit("client:disconnected", eventBytes)

	currentRound, err := h.repo.GetCurrentRoundByRoomID(ctx, roomID.String())
	if err != nil {
		return err
	}

	if currentRound.GuesserID == user.ID {
		room, err := h.repo.GetRoom(ctx, roomID)
		if err != nil {
			return err
		}

		remainingMembers := make([]uuid.UUID, 0, len(room.Members))
		for _, member := range room.Members {
			if member.UserID != user.ID {
				memberID, err := uuid.Parse(member.UserID)
				if err != nil {
					return err
				}
				remainingMembers = append(remainingMembers, memberID)
			}
		}

		if len(remainingMembers) == 0 {
			return errors.New("no remaining members")
		}

		newGuesserID := remainingMembers[rand.Intn(len(remainingMembers))]

		newRoundID, err := h.repo.CreateRound(ctx, currentRound.GameID, currentRound.RoundIndex, newGuesserID.String(), currentRound.KanjiID)
		if err != nil {
			return err
		}

		if err := h.repo.UpdateGameCurrentRound(ctx, currentRound.GameID, newRoundID); err != nil {
			return err
		}

		remainingDrawers := make([]string, 0, len(room.Members))
		for _, member := range room.Members {
			if member.UserID != user.ID {
				remainingDrawers = append(remainingDrawers, member.UserID)
			}
		}

		if len(remainingDrawers) == 0 {
			return errors.New("no remaining drawers")
		}

		newDrawerID := remainingDrawers[rand.Intn(len(remainingDrawers))]
		turnID, err := h.repo.CreateTurn(ctx, newRoundID, 1, newDrawerID)
		if err != nil {
			return err
		}

		if err := h.repo.UpdateRoundCurrentTurn(ctx, newRoundID, turnID); err != nil {
			return err
		}

		var kanjiStr string
		kanjiID := currentRound.KanjiID
		kanjiStr = kanjiID.String()

		roundStartedEvent := api.RoundStartedEvent{
			GuesserId:  newGuesserID.String(),
			RoundIndex: int(currentRound.RoundIndex),
			Kanji:      &kanjiStr,
		}

		roundStartedBytes, err := json.Marshal(roundStartedEvent)
		if err != nil {
			return err
		}

		s.To(socket.Room(roomID.String())).Emit("round:started", roundStartedBytes)
	}

	return nil
}
