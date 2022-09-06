package main

import (
	"fmt"
	"mp4stream/internal/server"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("FFMPEG")

	//create router instance
	router := mux.NewRouter()
	server := server.NewServer(router)

	server.Listen()

}
