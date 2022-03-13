package main

import (
	"net/http"

	"github.com/7onn/croupier/app/croupier-api/routes"
	"github.com/7onn/croupier/app/croupier-api/webserver"
	"github.com/gorilla/mux"
)

func BuildPipeline(srv webserver.Server, r *mux.Router) {
	r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)
	r.HandleFunc("/", routes.Home(srv)).Methods(http.MethodGet)
	r.HandleFunc("/play", routes.Play(srv)).Methods(http.MethodGet)
}
