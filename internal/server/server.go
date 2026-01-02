package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（局域网内）
	},
}

// Message WebSocket 消息结构
type Message struct {
	Type      string `json:"type"`
	Payload   string `json:"payload,omitempty"`
	AutoEnter bool   `json:"autoEnter,omitempty"`
	Success   bool   `json:"success,omitempty"`
}

// Handler WebSocket 处理器
type Handler struct {
	OnTextReceived func(text string, autoEnter bool)
	OnConnect      func()
	OnDisconnect   func()

	connMu sync.Mutex
	conn   *websocket.Conn
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	h.connMu.Lock()
	h.conn = conn
	h.connMu.Unlock()

	defer func() {
		conn.Close()
		h.connMu.Lock()
		h.conn = nil
		h.connMu.Unlock()
		if h.OnDisconnect != nil {
			h.OnDisconnect()
		}
	}()

	log.Println("📱 手机已连接:", r.RemoteAddr)
	if h.OnConnect != nil {
		h.OnConnect()
	}

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println("连接断开:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "text":
			if msg.Payload != "" && h.OnTextReceived != nil {
				h.OnTextReceived(msg.Payload, msg.AutoEnter)
			}
			// 发送确认
			response := Message{Type: "ack", Success: true}
			conn.WriteJSON(response)

		case "ping":
			conn.WriteJSON(Message{Type: "pong"})
		}
	}
}

// IsConnected 检查是否有客户端连接
func (h *Handler) IsConnected() bool {
	h.connMu.Lock()
	defer h.connMu.Unlock()
	return h.conn != nil
}
