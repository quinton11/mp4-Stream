package server

import (
	"log"
	handler "mp4stream/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// using gorilla mux library
type Server struct {
	Srv    *http.Server
	Router *mux.Router
}

func NewServer(router *mux.Router) *Server {
	var Server Server
	Server.Router = router
	Server.Srv = &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3000",
	}
	return &Server
}

func (server *Server) Listen() {
	//set everything settable
	server.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	//set handlefunctions
	server.Router.HandleFunc("/", handler.Home).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", server.Router))
}
