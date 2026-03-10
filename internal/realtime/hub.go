package realtime

import (
	"log"
	"sync"

	"github.com/AshrafAaref21/go-ws/internal/models"
)

type Hub struct {
	Clients map[int64]map[*Client]struct{}
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[int64]map[*Client]struct{}),
	}
}

func (h *Hub) BroadcastToAll(event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, clients := range h.Clients {
		for client := range clients {
			client.SendEvent(event)
		}
	}
}

func (h *Hub) GetClients(userId int64) ([]*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, ok := h.Clients[userId]
	if !ok || len(conns) == 0 {
		return nil, false
	}

	clients := make([]*Client, 0, len(conns))
	for client := range conns {
		clients = append(clients, client)
	}
	return clients, true
}

func (h *Hub) SendEventToUserIds(userIds []int64, sendId int64, eventType EventType, payload map[string]any) {
	for _, userId := range userIds {
		h.mu.RLock()
		clients, ok := h.Clients[userId]
		h.mu.RUnlock()

		if !ok {
			continue
		}

		for client := range clients {
			if client.User.ID == sendId {
				continue
			}
			client.SendEvent(Event{
				EventType: eventType,
				Payload:   payload,
			})
		}
	}
}

func (h *Hub) RegisterClientConnection(client *Client) {
	h.mu.Lock()

	conns, ok := h.Clients[client.User.ID]
	if !ok {
		conns = make(map[*Client]struct{})
		h.Clients[client.User.ID] = conns
	}

	conns[client] = struct{}{}
	firstConn := len(conns) == 1
	h.mu.Unlock()

	if firstConn {
		h.BroadcastToAll(Event{
			EventType: EventUserOnline,
			Payload:   client.User.ToMap(),
		})
		go func() {
			privates, err := models.GetPrivatesForUser(client.User.ID)
			if err != nil {
				log.Printf("failed to get privates: %v", err)
				return
			}
			for _, private := range privates {
				msgs, err := models.GetUndeliveredMessagesByPrivateID(private.ID)
				if err != nil {
					log.Printf("failed to get undelivered messages: %v", err)
					continue
				}
				for _, msg := range msgs {
					if msg.FromID == client.User.ID {
						continue
					}
					h.SendEventToUserIds([]int64{msg.FromID}, client.User.ID, EventDelivered, map[string]any{
						"message_id": msg.ID,
						"to_id":      client.User.ID,
					})
				}
			}
		}()
	}
}

func (h *Hub) UnregisterClientConnection(client *Client) {
	h.mu.Lock()
	conns, ok := h.Clients[client.User.ID]
	if ok {
		delete(conns, client)
		if len(conns) == 0 {
			delete(h.Clients, client.User.ID)
		}
	}
	lastConn := !ok || len(conns) == 0
	h.mu.Unlock()

	if lastConn {
		h.BroadcastToAll(Event{
			EventType: EventUserOffline,
			Payload:   client.User.ToMap(),
		})
	}
}

func (h *Hub) SendCurrentClients(toClient *Client) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]map[string]any, 0, len(h.Clients))
	seen := make(map[int64]struct{})

	for userId, conns := range h.Clients {
		if userId == toClient.User.ID {
			continue
		}
		_, ok := seen[userId]
		if ok {
			continue
		}

		for client := range conns {
			users = append(users, client.User.ToMap())
			seen[userId] = struct{}{}
			break
		}
	}
	toClient.Send <- Event{
		EventType: EventCurrentUsers,
		Payload:   users,
	}
}

func (h *Hub) SendError(clientId int64, message string) {
	h.mu.RLock()
	clients, ok := h.GetClients(clientId)
	h.mu.RUnlock()

	if !ok || len(clients) == 0 {
		return
	}

	for _, client := range clients {
		client.SendEvent(Event{
			EventType: EventError,
			Payload: map[string]any{
				"message": message,
			},
		})
	}
}

func (h *Hub) Shutdown() {
	h.mu.RLock()
	defer h.mu.RUnlock()

	log.Println("Shutting down Hub, notifying all clients.")

	for _, clients := range h.Clients {
		for client := range clients {
			client.SendEvent(Event{
				EventType: EventServerShutdown,
				Payload:   "Server is shutting down.",
			})
			client.Close()
		}
	}

	h.Clients = make(map[int64]map[*Client]struct{})
	log.Println("Hub shutdown complete.")
}
