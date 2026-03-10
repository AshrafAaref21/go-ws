package realtime

import (
	"sync"

	"github.com/AshrafAaref21/go-ws/internal/models"
	"github.com/coder/websocket"
)

type Client struct {
	User *models.User    `json:"user"`
	conn *websocket.Conn `json:"-"`
	Send chan Event      `json:"-"`
	once sync.Once       `json:"-"`
}

func NewClient(user *models.User, conn *websocket.Conn) *Client {
	return &Client{
		User: user,
		conn: conn,
		Send: make(chan Event, 512),
	}
}

func (c *Client) SendEvent(event Event) {
	select {
	case c.Send <- event:
	default:
	}
}

func (c *Client) Close() {
	c.once.Do(func() {
		if c.conn != nil {
			c.conn.Close(websocket.StatusNormalClosure, "Connection closed")
		}
		close(c.Send)
	})
}
