package main

import (
	"net/http"

	"github.com/7onn/croupier/routes"
	"github.com/7onn/croupier/storage"
	"github.com/7onn/croupier/webserver"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func init() {
	storage.Connect()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	webserver.New(webserver.Config{
		Port: 8080,
	}).Start(func(srv webserver.Server, r *mux.Router) {
		r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)
		r.HandleFunc("/", routes.Home(srv)).Methods(http.MethodGet)
		r.HandleFunc("/lobbyws", routes.LobbyWS(srv)).Methods(http.MethodGet)
		r.HandleFunc("/lobby", routes.Lobby(srv)).Methods(http.MethodGet)
		r.HandleFunc("/play", routes.Play(srv)).Methods(http.MethodGet)
	})
}
