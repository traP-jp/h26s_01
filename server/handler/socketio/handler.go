package socketio

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/traP-jp/h26s_01/server/model"
	"github.com/traP-jp/h26s_01/server/repository"
	"github.com/zishang520/socket.io/servers/socket/v3"
)

type Handler struct {
	ServeHandler http.Handler
	io           *socket.Server
	repo         *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	ioServer := socket.NewServer(nil, nil)

	return &Handler{
		ServeHandler: ioServer.ServeHandler(nil),
		io:           ioServer,
		repo:         repo,
	}
}

func (h *Handler) Start() (err error) {
	h.io.Use(func(soc *socket.Socket, next func(*socket.ExtendedError)) {
		// Directly access the underlying HTTP request headers
		userID := soc.Client().Request().Request().Header.Get("X-Forwarded-User")

		if userID == "" {
			next(socket.NewExtendedError("unauthorized: X-Forwarded-User header not found", nil))
			return
		}

		user, err := h.repo.GetOrCreateUser(context.Background(), userID)
		if err != nil {
			next(socket.NewExtendedError("failed to get or create user", err.Error()))
			return
		}

		soc.SetData(user)
		next(nil)
	})

	return h.io.On("connection", func(args ...any) {
		socket, ok := args[0].(*socket.Socket)

		if !ok {
			err = errors.New("socket assertion failed")

			return
		}

		h.registerEventHandlers(socket)
	})
}

func (h *Handler) getLoggedInUser(socket *socket.Socket) (*model.User, error) {
	// Directly access the underlying HTTP request headers
	userID := socket.Client().Request().Request().Header.Get("X-Forwarded-User")

	return h.repo.GetOrCreateUser(context.Background(), userID)
}

func createEventListenerForHandlersWithoutBody(socket *socket.Socket, handler func(socket *socket.Socket) error) func(args ...any) {
	return func(args ...any) {
		if err := handler(socket); err != nil {
			slog.Error("handling event", "error", err)
		}
	}
}

func createEventListenerForHandlersWithBody[T any](socket *socket.Socket, handler func(socket *socket.Socket, event T) error) func(args ...any) {
	return func(args ...any) {
		var bodyBytes []byte
		var err error

		switch arg := args[0].(type) {
		case []byte:
			bodyBytes = arg
		case string:
			bodyBytes = []byte(arg)
		default:
			bodyBytes, err = json.Marshal(arg)
			if err != nil {
				slog.Error("marshaling event body", "error", err)

				return
			}
		}

		var event T

		if err := json.Unmarshal(bodyBytes, &event); err != nil {
			slog.Error("unmarshaling event", "error", err)

			return
		}

		if err := handler(socket, event); err != nil {
			slog.Error("handling event", "error", err)
		}
	}
}

func (h *Handler) registerEventHandlers(socket *socket.Socket) {
	slog.Info("Registering event handlers", "socketID", socket.Id())
	socket.On("room:join", createEventListenerForHandlersWithBody(socket, h.handleJoinRoom))
	socket.On("game:ready", createEventListenerForHandlersWithoutBody(socket, h.handleGameReady))
	socket.On("draw:stroke", createEventListenerForHandlersWithBody(socket, h.handleDrawStroke))
	socket.On("answer:submit", createEventListenerForHandlersWithBody(socket, h.handleAnswerSubmit))
	socket.On("round:end", createEventListenerForHandlersWithoutBody(socket, h.handleRoundEnd))
	socket.On("client:disconnect", createEventListenerForHandlersWithoutBody(socket, h.handleClientDisconnect))

	go func() {
		if err := h.handleRoomListUpdated(socket); err != nil {
			slog.Error("handling room list updated", "error", err)
		}
	}()
	go func() {
		if err := h.roomUpdatedEventHandler(socket); err != nil {
			slog.Error("handling room updated", "error", err)
		}
	}()
	go func() {
		if err := h.handleRoundAnswer(socket); err != nil {
			slog.Error("handling round answer", "error", err)
		}
	}()
	go func() {
		if err := h.handleRoundStartedEvent(socket); err != nil {
			slog.Error("handling round started", "error", err)
		}
	}()
	go func() {
		if err := h.handleTurnStartedEvent(socket); err != nil {
			slog.Error("handling turn started", "error", err)
		}
	}()
}
