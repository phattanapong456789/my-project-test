package handlers

import (
	"auth-app/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type reservationsHub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func newReservationsHub() *reservationsHub {
	return &reservationsHub{clients: map[*websocket.Conn]struct{}{}}
}

func (h *reservationsHub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = struct{}{}
}

func (h *reservationsHub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *reservationsHub) broadcast(v any) {
	h.mu.Lock()
	clients := make([]*websocket.Conn, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.Unlock()

	for _, c := range clients {
		_ = c.SetWriteDeadline(time.Now().Add(3 * time.Second))
		if err := c.WriteJSON(v); err != nil {
			_ = c.Close()
			h.remove(c)
		}
	}
}

var globalReservationsHub = newReservationsHub()

func BroadcastReservationsChanged() {
	globalReservationsHub.broadcast(gin.H{"type": "reservations_changed", "ts": time.Now().Unix()})
}

func (h *TableAdminHandler) ReservationsWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}
	if claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: admin only"})
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	globalReservationsHub.add(conn)

	// Send initial ping/event so client knows it's connected
	_ = conn.WriteJSON(gin.H{"type": "connected"})

	// Read loop: keep connection alive; if client closes => exit
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	globalReservationsHub.remove(conn)
	_ = conn.Close()
}
