package main

import (
	"net/http"

	"github.com/7onn/croupier/routes"
	"github.com/7onn/croupier/webserver"
	"github.com/gorilla/mux"
)

func init() {}

func main() {
	webserver.New(webserver.Config{
		Port: 8080,
	}).Start(func(srv webserver.Server, r *mux.Router) {
		r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)
		r.HandleFunc("/", routes.Home(srv)).Methods(http.MethodGet)
		r.HandleFunc("/play", routes.Play(srv)).Methods(http.MethodGet)
	})
}
