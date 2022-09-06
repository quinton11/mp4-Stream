package server

import (
	"log"
	handler "mp4stream/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Srv     *http.Server
	Router  *mux.Router
	Handler *handler.Handler
}

func NewServer(router *mux.Router) *Server {
	var Server Server
	Server.Router = router
	Server.Handler = handler.NewHandler()
	return &Server
}

func (server *Server) Listen() {
	//Serve statics
	server.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//set handlefunctions
	server.Router.HandleFunc("/", server.Handler.Home).Methods("GET")
	server.Router.HandleFunc("/signal", server.Handler.Signal).Methods("POST")
	server.Router.HandleFunc("/stream", server.Handler.Stream).Methods("POST")

	//Initialize webRtc Agent
	server.Handler.Agent.InitProcess()

	log.Fatal(http.ListenAndServe(":3000", server.Router))
}
