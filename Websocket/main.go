package Websocket

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"plc-backend/Shm"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleNewConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return ws, err
}

func writer(conn *websocket.Conn) {
	for {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			data, err := Shm.Read(os.Getenv("SHM-R-DATA"), os.Getenv("SHM-R-LOCK"))

			if err != nil {
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
				if errors.Is(err, syscall.EPIPE) {
					return
				}
				fmt.Println(err)
				return
			}
		}
	}
}

func Write(w http.ResponseWriter, r *http.Request) {
	ws, err := handleNewConnection(w, r)

	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	go writer(ws)
}
