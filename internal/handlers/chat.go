package handlers

import (
	"net/http"
	"sync"

	"SharkLava/random_chat/pkg/queue"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	send    chan []byte
	partner *Client
}

type Hub struct {
	userQueue *queue.Queue
	mu        sync.Mutex
}

func NewHub(userQueue *queue.Queue) *Hub {
	return &Hub{
		userQueue: userQueue,
	}
}

func (h *Hub) Run() {
	for {
		if h.userQueue.Len() >= 2 {
			h.mu.Lock()
			user1 := h.userQueue.Pop().(*Client)
			user2 := h.userQueue.Pop().(*Client)
			h.mu.Unlock()

			go h.pairClients(user1, user2)
		}
	}
}

func (h *Hub) pairClients(client1, client2 *Client) {
	client1.partner = client2
	client2.partner = client1

	client1.send <- []byte("You've been paired with a stranger. Start chatting!")
	client2.send <- []byte("You've been paired with a stranger. Start chatting!")

	go client1.readPump()
	go client2.readPump()
}

func ServeHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "templates/index.html")
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

	hub.mu.Lock()
	hub.userQueue.Push(client)
	hub.mu.Unlock()

	// Send waiting message
	client.send <- []byte("Waiting for a partner...")

	go client.writePump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.mu.Lock()
		if c.partner != nil {
			c.partner.send <- []byte("Your partner has disconnected.")
			c.partner.partner = nil
		}
		c.hub.userQueue.Remove(c)
		c.hub.mu.Unlock()
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		if c.partner != nil {
			c.partner.send <- message
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				return
			}
		}
	}
}
