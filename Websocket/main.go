package Websocket

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"plc-backend/Shm"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var sockets = make(map[string]*Connection)

// Connection Websocket connection interface
type Connection struct {
	Socket *websocket.Conn
	sync.Mutex
}

// Send Method to send data to connection
func (c *Connection) Send(message string) error {
	c.Lock()
	defer c.Unlock()
	return c.Socket.WriteMessage(websocket.TextMessage, []byte(message))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Function to handle new websocket connection
func handleNewConnection(w http.ResponseWriter, r *http.Request) error {
	// Upgrade the http connection
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create unique id for the connection
	id := uuid.New().String()
	// Create a new connection object
	conn := &Connection{Socket: ws}
	// Register the connection
	sockets[id] = conn

	return nil
}

// Function that writes to all registered connections
func write() {
	for {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			data, err := Shm.Read(os.Getenv("SHM-R-DATA"), os.Getenv("SHM-R-LOCK"))

			if err != nil {
				return
			}

			for key, connection := range sockets {
				if err := connection.Send(data); err != nil {
					if errors.Is(err, syscall.EPIPE) {
						return
					}
					delete(sockets, key)
				}
			}
		}
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Upgrade and register connection
	err := handleNewConnection(w, r)
	// Check for error in connection
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}
	// Write to all registered connections
	go write()
}
