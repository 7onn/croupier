package routes

import (
	"net/http"

	"github.com/7onn/croupier/app/croupier-api/webserver"
	ws "github.com/7onn/croupier/pkg/websocket"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Play(srv webserver.Server) http.HandlerFunc {
	hub := ws.NewHub()
	go hub.Run()
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err)
			return
		}
		client := &ws.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		go client.WritePump()
		go srv.PlayPoker(client)
	}
}
