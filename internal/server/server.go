package server

import (
	"log"
	handler "mp4stream/internal/handlers"
	"net/http"

	"github.com/fatih/color"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

type Server struct {
	Srv     *http.Server
	Router  *mux.Router
	Handler *handler.Handler
}

func NewServer(router *mux.Router, path string) *Server {
	var Server Server
	Server.Router = router
	Server.Handler = handler.NewHandler(path)
	return &Server
}

func (server *Server) Listen() {
	//
	port := ":3000"
	colord := color.New(color.FgHiBlue, color.Bold).Add(color.BgHiGreen)

	//Serve statics
	server.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	//set handlefunctions
	//server.Router.HandleFunc("/", server.Handler.Home).Methods("GET")
	server.Router.HandleFunc("/signal", server.Handler.Signal).Methods("POST")
	server.Router.HandleFunc("/streamup", server.Handler.StreamUp).Methods("POST")

	//Initialize webRtc Agent
	colord.Println("Initiating Agent...")
	server.Handler.Agent.InitProcess()
	url := "http://127.0.0.1" + port
	colord.Println("Serving at:", url)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowCredentials: true,
	})
	handler := c.Handler(server.Router)
	log.Fatal(http.ListenAndServe(port, handler))

}
