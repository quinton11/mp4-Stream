package main

import (
	"mp4stream/internal/server"

	"github.com/gorilla/mux"
)

func main() {

	//create router instance
	router := mux.NewRouter()
	server := server.NewServer(router, "")

	server.Listen()

}
