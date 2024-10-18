package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/watermill"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust this for production use
	},
}

type WSMessage struct {
	Event string      `json:"event"`
	Body  interface{} `json:"body"`
}

type WebSocketServer struct {
	clients      map[string][]*websocket.Conn // Changed to map user ID to a slice of connections
	clientsMutex sync.RWMutex
	tracer       trace.Tracer
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[string][]*websocket.Conn),
		tracer:  otel.Tracer("websocket-server"),
	}
}

func (s *WebSocketServer) HandleWebSocket(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		propagator := otel.GetTextMapPropagator()
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		ctx, span := s.tracer.Start(ctx, "websocket_connection", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		span.SetAttributes(
			attribute.String("client_ip", c.ClientIP()),
			attribute.String("user_agent", c.Request.UserAgent()),
		)

		userID := c.Query("user_id")
		if userID == "" {
			logger.AppLogger.Error(ctx, "Missing user ID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
			span.RecordError(errors.New("missing user id"))
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to upgrade connection",
				zap.Error(err),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_id", userID),
			)
			span.RecordError(err)
			return
		}

		logger.AppLogger.Info(ctx, "Client connected successfully",
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_id", userID),
		)

		s.clientsMutex.Lock()
		s.clients[userID] = append(s.clients[userID], conn)
		s.clientsMutex.Unlock()

		defer func() {
			s.removeConnection(userID, conn)
			conn.Close()
		}()

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				logger.AppLogger.Error(ctx, "Error reading message",
					zap.Error(err),
					zap.String("client_ip", c.ClientIP()),
					zap.String("user_id", userID),
					zap.ByteString("payload", p),
				)
				span.RecordError(err)
				s.sendErrorMessage(ctx, conn, "Error reading message")
				break
			}

			var wsmessage WSMessage
			if err = json.Unmarshal(p, &wsmessage); err != nil {
				logger.AppLogger.Error(ctx, "Failed to parse ws msg",
					zap.Int("ms type", messageType),
					zap.ByteString("payload", p),
					zap.String("user_id", userID),
				)
				s.sendErrorMessage(ctx, conn, "Failed to parse message")
				continue
			}

			if wsmessage.Event == "" {
				logger.AppLogger.Error(ctx, "Missing event",
					zap.Int("ms type", messageType),
					zap.ByteString("payload", p),
					zap.String("user_id", userID),
				)
				s.sendErrorMessage(ctx, conn, "Missing event in message")
				continue
			}

			logger.AppLogger.Info(ctx, "Message received",
				zap.Int("ms type", messageType),
				zap.ByteString("payload", p),
				zap.String("user_id", userID),
			)

			_, msgSpan := s.tracer.Start(ctx, "websocket_message", trace.WithSpanKind(trace.SpanKindServer))
			msgSpan.SetAttributes(attribute.Int("message.type", messageType))

			if err := watermill.PublishReceivedWSMsgEvent(
				ctx,
				appCtx.GetLocalPubSub().GetUnblockPubSub(),
				p,
			); err != nil {
				logger.AppLogger.Error(ctx, "Error writing message",
					zap.Error(err),
					zap.String("client_ip", c.ClientIP()),
					zap.ByteString("payload", p),
					zap.String("user_id", userID),
				)
			}

			msgSpan.End()
		}
	}
}

func (s *WebSocketServer) removeConnection(userID string, conn *websocket.Conn) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	connections := s.clients[userID]
	for i, c := range connections {
		if c == conn {
			s.clients[userID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	if len(s.clients[userID]) == 0 {
		delete(s.clients, userID)
	}
}

func (s *WebSocketServer) sendErrorMessage(ctx context.Context, conn *websocket.Conn, errorMessage string) {
	errorMsg := WSMessage{
		Event: "error_message",
		Body:  errorMessage,
	}
	jsonMsg, err := json.Marshal(errorMsg)
	if err != nil {
		logger.AppLogger.Error(ctx, "Error marshaling error message",
			zap.Error(err),
		)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
		logger.AppLogger.Error(ctx, "Error sending error message",
			zap.Error(err),
		)
	}
}

func (s *WebSocketServer) BroadcastMessage(ctx context.Context, message []byte) {
	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	for userID, connections := range s.clients {
		for _, conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				logger.AppLogger.Error(ctx, "Error broadcasting message",
					zap.Error(err),
					zap.Int("message_size", len(message)),
					zap.String("user_id", userID),
				)
				s.removeConnection(userID, conn)
			}
		}
	}
}

func (s *WebSocketServer) SendMessageToUser(ctx context.Context, userID string, message []byte) error {

	s.clientsMutex.RLock()
	connections, exists := s.clients[userID]
	s.clientsMutex.RUnlock()

	if !exists || len(connections) == 0 {
		err := errors.New("user not found or has no active connections")
		logger.AppLogger.Error(ctx, "User not found or has no active connections",
			zap.String("user_id", userID),
		)
		return err
	}

	for _, conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logger.AppLogger.Error(ctx, "Error sending message to user",
				zap.Error(err),
				zap.String("user_id", userID),
				zap.Int("message_size", len(message)),
			)
			s.removeConnection(userID, conn)
		} else {
			logger.AppLogger.Info(ctx, "Message sent to user successfully",
				zap.String("user_id", userID),
				zap.Int("message_size", len(message)),
			)
		}
	}

	return nil
}
