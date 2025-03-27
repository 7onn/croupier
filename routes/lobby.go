package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/7onn/croupier/middleware"
	"github.com/7onn/croupier/storage"
	"github.com/7onn/croupier/templates"
	"github.com/7onn/croupier/webserver"
	ws "github.com/7onn/croupier/websocket"
	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Lobby is for
func Lobby(srv webserver.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("jwt")
		table, err := storage.GetFormingTable()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cookie == nil {
			templ.Handler(templates.Lobby(table, nil)).ServeHTTP(w, r)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return middleware.JWTKey(), nil
		})
		if err != nil {
			w.Header().Set("Set-Cookie", "jwt=")
			templ.Handler(templates.Lobby(table, nil)).ServeHTTP(w, r)
			return
		}
		if token.Valid {
			user := storage.GetUserByID(token.Claims.(jwt.MapClaims)["user_id"].(string))
			if user != nil {
				// table, err := storage.GetTableByID(token.Claims.(jwt.MapClaims)["table_id"].(string))
				// if err != nil {
				// 	log.Err(err)
				// }
				templ.Handler(templates.Lobby(table, user)).ServeHTTP(w, r)
				return
			}
		}
	}
}

var mutex sync.Mutex

func LobbyWS(srv webserver.Server) http.HandlerFunc {
	hub := ws.NewHub()
	go hub.Run()
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err)
			return
		}
		client := &ws.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
		client.Hub.Register <- client

		go LobbyWsWrite(client)
		go LobbyWsRead(client)
		mutex.Unlock()
	}
}

func LobbyWsRead(c *ws.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(ws.MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(ws.PongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, ws.Newline, ws.Space, -1))

		handleClientMessage(c, message)
	}
}

func handleClientMessage(c *ws.Client, message []byte) {
	fmt.Printf("client sent: %s\n", string(message))

	type Message struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	var baseMsg Message
	if err := json.Unmarshal(message, &baseMsg); err != nil {
		log.Err(err)
		return
	}

	switch baseMsg.Type {
	case "members":
		type MembersMessage struct {
			TableID uuid.UUID `json:"tableId"`
		}
		var membersMsg *MembersMessage
		err := json.Unmarshal(baseMsg.Data, &membersMsg)
		if err != nil {
			log.Err(err)
			return
		}

		table, err := storage.GetTableByID(membersMsg.TableID.String())
		if err != nil {
			log.Err(err)
			return
		}

		members, _ := json.Marshal(table.Members)
		msg := Message{
			Type: "members",
			Data: members,
		}
		jsonMsg, _ := json.Marshal(msg)
		c.Hub.Broadcast <- jsonMsg

	case "newJoiner":
		type NewJoinerMessage struct {
			Name    string `json:"name"`
			TableID string `json:"tableId"`
		}
		var newJoinerMsg *NewJoinerMessage
		if err := json.Unmarshal(baseMsg.Data, &newJoinerMsg); err != nil {
			log.Err(err)
			return
		}
		fmt.Printf("new joiner: %s\n", newJoinerMsg.Name)
		user, err := storage.CreateUser(newJoinerMsg.Name, newJoinerMsg.TableID)
		if err != nil {
			log.Err(err)
			return
		}

		table, err := storage.GetTableByID(newJoinerMsg.TableID)
		if err != nil {
			log.Err(err)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  user.ID.String(),
			"table_id": table.ID.String(),
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString(middleware.JWTKey())
		if err != nil {
			log.Err(err)
			return
		}

		members, _ := json.Marshal(table.Members)
		msg := Message{
			Type: "members",
			Data: members,
		}
		jsonMsg, _ := json.Marshal(msg)
		c.Hub.Broadcast <- jsonMsg

		msg = Message{
			Type: "jwt",
			Data: []byte(fmt.Sprintf(`"%s"`, tokenString)),
		}
		jsonMsg, _ = json.Marshal(msg)
		c.Send <- jsonMsg
		return

	default:
		fmt.Println("Unknown message type:", baseMsg.Type)
	}

}

func LobbyWsWrite(c *ws.Client) {
	ticker := time.NewTicker(ws.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetReadLimit(ws.MaxMessageSize)
			c.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// send to client
			w.Write(message)
			fmt.Printf("server sent: %s\n", string(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(ws.Newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
