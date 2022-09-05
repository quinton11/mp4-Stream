package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"mp4stream/internal/service"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var newPc = true

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

func (h *Handler) GetCmd(input string) string {
	inputArr := strings.Split(input, " ")
	return inputArr[0]
}

func (h *Handler) Parse(input []byte) (map[string]interface{}, bool) {
	JSON := make(map[string]interface{})
	err := json.Unmarshal(input, &JSON)
	if err != nil {
		fmt.Println(err)
	}

	_, okOffer := JSON["type"]
	if okOffer {
		return JSON, okOffer
	}
	return JSON, false
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	//create websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	h.Agent.Ws = conn
	if err != nil {
		panic(err)
	}

	//continuosly read in messages
	go func() {
		for {

			_, mes, err := h.Agent.Ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
			}
			//cmd := h.GetCmd(string(mes))
			//fmt.Println(cmd)
			if string(mes) != "undefined" {
				offer := webrtc.SessionDescription{}
				remoteCandidate := webrtc.ICECandidate{}
				res, isOffer := h.Parse(mes)
				//fmt.Println(res)
				//If message is offer
				if isOffer {
					json.Unmarshal(mes, &offer)
					fmt.Println(offer)
					if h.Agent.Pconnect.CurrentRemoteDescription() == &offer {
						fmt.Println("Same offer")
					}
					//set remote description with offer
					if errR := h.Agent.Pconnect.SetRemoteDescription(offer); errR != nil {
						fmt.Println("Remote Description")
						panic(errR)
					}
					//create answer
					answer, err := h.Agent.Pconnect.CreateAnswer(nil)
					if err != nil {
						fmt.Println("Creating Answer")
						panic(err)
					}
					//set local description with answer
					if errA := h.Agent.Pconnect.SetLocalDescription(answer); errA != nil {
						fmt.Println("Local Description")
						panic(errA)
					}

					//send response to client

					fmt.Println("Writing")
					if errWr := h.Agent.Ws.WriteJSON(answer); errWr != nil {
						panic(errWr)
					}

					if newPc {
						fmt.Println("Peer connection registered")
						newPc = !newPc
						fmt.Printf("\nSwitching states: %v", newPc)
					}
				}

				//If message is ICE candidate
				if !isOffer {
					fmt.Println(res)
					json.Unmarshal(mes, &remoteCandidate)
					fmt.Println(remoteCandidate)
					fmt.Println(res)
				}

			}

		}
	}()
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
	//gcomplete := webrtc.GatheringCompletePromise(h.Agent.Pconnect)
	//Set answer as local SDP description
	h.Agent.Pconnect.SetLocalDescription(answer)
	fmt.Println(answer)

	ans, err := json.Marshal(answer)
	if err != nil {
		fmt.Println(err)
	}

	/* //starting stream
	go func() {
		//<-gcomplete
		//h.Agent.StreamTrack() //push ffmpeg buffers unto localtrack
		h.Agent.StreamRTP() //push ffmpeg stream unto RTP localtrack
	}() */

	w.WriteHeader(http.StatusOK)
	w.Write(ans)
	//fmt.Fprintf(w, "Signalling, remote answer...")
}
