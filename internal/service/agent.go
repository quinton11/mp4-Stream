package service

import (
	//"bytes"
	//"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	//ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

// Agent handles peerconnection.
// for creation of offer and answer
// for adding web tracks
type Agent struct {
	Pconnect *webrtc.PeerConnection
	Track    *webrtc.TrackLocalStaticSample
	Ws       *websocket.Conn
}

type Offer struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}

//Create new PeerConnection  Agent

func NewAgent() (*Agent, error) {
	//webrtc server configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.1.google.com:19302", "stun:stun.2.google.com:19302"},
			},
		},
	}
	peerconnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return &Agent{Pconnect: peerconnection}, nil
}

func (agent *Agent) InitProcess() (*webrtc.TrackLocalStaticSample, error) {
	//this creates local track
	localtrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264"}, "video", "pion1")
	if err != nil {
		return nil, err
	}

	//add track to peerconnection
	agent.Pconnect.AddTrack(localtrack)
	agent.Track = localtrack

	//create onice change listener
	agent.Pconnect.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("New Connection state %v", connectionState)
	})

	agent.Pconnect.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		fmt.Printf("Ice Candidate %v \n", candidate)
		if candidate != nil {
			//agent.Pconnect.AddICECandidate(candidate.ToJSON())
			//On every Ice candidate received, send the updated
			//local description
			agent.Ws.WriteJSON(agent.Pconnect.CurrentLocalDescription())
			//go agent.Ws.ReadJSON(agent.Pconnect.RemoteDescription())

		}
	})

	agent.Pconnect.OnNegotiationNeeded(func() {
		fmt.Println("Negotiation Needed")
	})

	//return track
	return localtrack, nil
}

func (agent *Agent) CreateOffer() (*webrtc.SessionDescription, error) {
	offer, err := agent.Pconnect.CreateOffer(&webrtc.OfferOptions{ICERestart: true})
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

func (agent *Agent) StreamTrack() {
	//Load movie file
	//get working dir
	fmt.Println("Streaming...")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := "\\assets\\opbest.mp4"
	//check if file exists
	if _, err := os.Stat(dir + filename); err != nil {
		fmt.Println("File missing")
		fmt.Println(err)
		//NB: Write error handler to handle case
		//for missing file
	}

	//start ffmpeg process
	//Tell ffmpeg to read from file into stdout.
	//Then we read from stdout and write to track
	//args

	//using exec
	command := "ffmpeg"
	args := []string{
		"-i",
		dir + filename,
		"-c:v",
		"libx264",
		"-preset",
		"superfast",
		"-f",
		"ismv",
		"pipe:1",
	} //output to stdout
	cmd := exec.Command(command, args...)
	stdout, errP := cmd.StdoutPipe()
	if errP != nil {
		panic(errP)
	}

	errS := cmd.Start()
	if errS != nil {
		panic(errS)
	}

	fmt.Println("FFMPEG started...")

	//set stdout as output to ffmpeg_go .run()
	//create reader to read from stdout by setting os.stdout as input to
	//new reader.
	//in another go routine, read outputs to stdout into a buffer with a bitrate
	//write bytes from buffer into localstatictrack

	//addtrack

	//works
	fmt.Println("Reading from STDOUT...")
	buf := make([]byte, 1024*64)
	for {
		//Reading from stdout
		n, err := stdout.Read(buf)
		//fmt.Println(buf[:n])
		if err != nil {
			fmt.Println(err)

			if err == io.EOF {
				fmt.Println("Done.")
				break
			}
		}

		//write to samplet
		err = agent.Track.WriteSample(media.Sample{Data: buf[:n]})
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

}
