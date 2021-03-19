package controllers

import (
	"croupier/pkg/user"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

//Play !
var Play = func(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	j := q.Get("jwt")
	room := q.Get("room")
	fmt.Printf("\n jwt %+v \n room: %+v", j, room)

	tk := &user.Token{}
	jwt.ParseWithClaims(j, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN")), nil
	})

	fmt.Println(tk.UserID)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
