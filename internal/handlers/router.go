package handler

import (
	"log"
	"mp4stream/internal/service"
	"net/http"
)

type Handler struct {
	Agent *service.Agent
}

func NewHandler() *Handler {
	agent, err := service.NewAgent()
	if err != nil {
		log.Fatal(err)
	}
	return &Handler{Agent: agent}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	//http.ServeFile(w, r, http.FileServer(http.Dir("./static")))

	w.Write([]byte("Welcome👋🏾"))
}

func (h *Handler) Signal(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//b := make([]byte, 100)
	//fmt.Println(r.Body.Read(b))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Signalling, remote answer..."))
	//fmt.Fprintf(w, "Signalling, remote answer...")
}
