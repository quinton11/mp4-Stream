package service

import (
	"fmt"

	"github.com/pion/webrtc/v3"
)

// Agent handles peerconnection.
// for creation of offer and answer
// for adding web tracks
type Agent struct {
	Pconnect *webrtc.PeerConnection
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

	//create onice change listener
	agent.Pconnect.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("New Connection state %v", connectionState)
	})

	//return track
	return localtrack, nil
}

func (agent *Agent) CreateOffer() (*webrtc.SessionDescription, error) {
	offer, err := agent.Pconnect.CreateOffer(&webrtc.OfferOptions{})
	if err != nil {
		return nil, err
	}
	return &offer, nil
}
