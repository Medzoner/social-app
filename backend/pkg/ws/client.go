package ws

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"social-app/internal/connector"
)

const (
	UserKeyPrefix   = "ws:user:"
	SocketKeyPrefix = "ws:socket:"
	ttlTimeout      = 15 * time.Minute
)

type Connector struct {
	Clients      map[uint64]*websocket.Conn
	redisConn    *connector.RedisConnector
	Upgrader     websocket.Upgrader
	ClientsMutex sync.RWMutex
}

func NewConnector(r *connector.RedisConnector) *Connector {
	return &Connector{
		Clients:      make(map[uint64]*websocket.Conn),
		ClientsMutex: sync.RWMutex{},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		redisConn: r,
	}
}

func (c *Connector) ConnectClient(r *http.Request, w http.ResponseWriter, userID uint64) (*websocket.Conn, string, error) {
	conn, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, "", fmt.Errorf("upgrade error: %w", err)
	}

	socketID := uuid.NewString()
	if err := c.ReplaceClient(userID, conn); err != nil {
		return nil, "", err
	}

	userKey := fmt.Sprintf("%s%d", UserKeyPrefix, userID)
	socketKey := fmt.Sprintf("%s%s", SocketKeyPrefix, socketID)

	if err := c.redisConn.Set(r.Context(), userKey, socketID, ttlTimeout); err != nil {
		return nil, "", fmt.Errorf("redis set user→socket error: %w", err)
	}
	if err := c.redisConn.Set(r.Context(), socketKey, fmt.Sprintf("%d", userID), ttlTimeout); err != nil {
		return nil, "", fmt.Errorf("redis set socket→user error: %w", err)
	}

	return conn, socketID, nil
}

func (c *Connector) CloseClient(ctx context.Context, userID uint64) error {
	id, _, err := c.redisConn.DeleteIfExists(ctx, fmt.Sprintf("%s%d", UserKeyPrefix, userID))
	if err != nil {
		return fmt.Errorf("redis delete error: %w", err)
	}
	if id != "" {
		_, _, err := c.redisConn.DeleteIfExists(ctx, fmt.Sprintf("%s%s", SocketKeyPrefix, id))
		if err != nil {
			return fmt.Errorf("redis delete error: %w", err)
		}
	}

	c.ClientsMutex.Lock()
	defer c.ClientsMutex.Unlock()

	conn, ok := c.Clients[userID]
	if !ok || conn == nil {
		log.Printf("No active connection for user %d", userID)
		return nil
	}
	delete(c.Clients, userID)

	if err := conn.Close(); err != nil {
		return fmt.Errorf("close error: %w", err)
	}

	return nil
}

func (c *Connector) closeAll(ctx context.Context) error {
	c.ClientsMutex.Lock()
	defer c.ClientsMutex.Unlock()

	for userID, conn := range c.Clients {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("Error closing connection for user %d: %v", userID, err)
			}
		}

		id, _, err := c.redisConn.DeleteIfExists(ctx, fmt.Sprintf("%s%d", UserKeyPrefix, userID))
		if err != nil {
			return fmt.Errorf("redis delete error: %w", err)
		}
		if id != "" {
			_, _, err := c.redisConn.DeleteIfExists(ctx, fmt.Sprintf("%s%s", SocketKeyPrefix, id))
			if err != nil {
				return fmt.Errorf("redis delete error: %w", err)
			}
		}
	}
	c.Clients = make(map[uint64]*websocket.Conn)
	return nil
}

func (c *Connector) Send(userID uint64, notif any) error {
	c.ClientsMutex.Lock()
	conn, ok := c.Clients[userID]
	c.ClientsMutex.Unlock()

	if !ok || conn == nil {
		return fmt.Errorf("no connection for user %d", userID)
	}

	if err := conn.WriteJSON(notif); err != nil {
		return fmt.Errorf("failed to write JSON to user %d: %w", userID, err)
	}
	return nil
}

func (c *Connector) ReplaceClient(userID uint64, newConn *websocket.Conn) error {
	c.ClientsMutex.Lock()
	defer c.ClientsMutex.Unlock()

	if oldConn, ok := c.Clients[userID]; ok && oldConn != nil {
		if err := oldConn.Close(); err != nil {
			return fmt.Errorf("failed to close old connection for user %d: %w", userID, err)
		}
	}

	c.Clients[userID] = newConn
	return nil
}

func (c *Connector) IsUserOnline(ctx context.Context, id uint64) (bool, error) {
	exists, err := c.redisConn.IsExists(ctx, fmt.Sprintf("%s%d", UserKeyPrefix, id))
	if err != nil {
		return false, fmt.Errorf("redis error: %w", err)
	}
	return exists, nil
}
