package http

import (
	"context"
	"homework/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocketHandler struct {
	useCases UseCases
	shutdown context.Context
	cancel   context.CancelFunc
}

func NewWebSocketHandler(useCases UseCases) *WebSocketHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &WebSocketHandler{
		useCases: useCases,
		shutdown: ctx,
		cancel:   cancel,
	}
}

func (h *WebSocketHandler) Handle(ctx *gin.Context, id int64) error {
	if _, err := h.useCases.Sensor.GetSensorByID(ctx, id); err != nil {
		return err
	}

	conn, err := websocket.Accept(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	defer conn.Close(websocket.StatusNormalClosure, "bye-bye")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	closeCtx := conn.CloseRead(h.shutdown)

	var event *domain.Event

	for {
		select {
		case <-closeCtx.Done():
			return closeCtx.Err()
		case <-ticker.C:
			newEvent, err := h.useCases.Event.GetLastEventBySensorID(ctx, id)
			if err != nil {
				continue
			}

			if event != nil && newEvent.Timestamp.Equal(event.Timestamp) {
				continue
			}

			event = newEvent

			if err := wsjson.Write(ctx, conn, event); err != nil {
				return err
			}
		}
	}
}

func (h *WebSocketHandler) Shutdown() error {
	h.cancel()
	time.Sleep(time.Second)
	return nil
}
