package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mp4stream/internal/service"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	Agent *service.Agent
	Ws    *websocket.Conn
}

func NewHandler() *Handler {
	agent, err := service.NewAgent()
	if err != nil {
		log.Fatal(err)
	}
	return &Handler{Agent: agent}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	//create websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	h.Agent.Ws = conn
	if err != nil {
		panic(err)
	}
	//w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	//http.ServeFile(w, r, http.FileServer(http.Dir("./static")))

	w.Write([]byte("Welcome👋🏾"))
}

func (h *Handler) Signal(w http.ResponseWriter, r *http.Request) {
	//Read in stream of json data and store in placeholder
	JSON := make(map[string]interface{})
	d := json.NewDecoder(r.Body)
	//d.DisallowUnknownFields()
	d.Decode(&JSON)
	//print out data
	fmt.Println(JSON)
	b, err := json.Marshal(JSON)
	if err != nil {
		fmt.Println("Unmarshalling error")
		fmt.Println(err)
	}

	var offer webrtc.SessionDescription
	json.Unmarshal(b, &offer)

	//Set remote SDP description
	h.Agent.Pconnect.SetRemoteDescription(offer)
	//Create Answer to Offer
	answer, err := h.Agent.Pconnect.CreateAnswer(nil)
	if err != nil {
		fmt.Println(err)
	}
	//gather ICE candidates
	gcomplete := webrtc.GatheringCompletePromise(h.Agent.Pconnect)
	//Set answer as local SDP description
	h.Agent.Pconnect.SetLocalDescription(answer)
	fmt.Println(answer)

	ans, err := json.Marshal(answer)
	if err != nil {
		fmt.Println(err)
	}

	//starting stream
	go func() {
		<-gcomplete
		h.Agent.StreamTrack() //push ffmpeg buffers unto localtrack
	}()

	w.WriteHeader(http.StatusOK)
	w.Write(ans)
	//fmt.Fprintf(w, "Signalling, remote answer...")
}
