package main

import (
	"fmt"
	"log"

	"github.com/pion/webrtc/v3"
)

func main() {
	fmt.Println("Hello Stream")

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.1.google.com:19302", "stun:stun.2.google.com:19302"},
			},
		},
	}

	//create peer connection
	peerconnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	//create webrtc signal description
	localTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/H.264"}, "video", "pion1")
	if err != nil {
		log.Fatal(err)
	}
	peerconnection.AddTrack(localTrack)
	//create local track - its transmitting streams to browser
	//so no need for remote tracks
	peerconnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("New Connection state %v", connectionState)
	})
	//receive offer from remote peer and establish remote SDP
	offer, err := peerconnection.CreateOffer(&webrtc.OfferOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Offer: %v", offer)
	//create answer and establish local SDP and send response to
	//remote peer
	//gst.CreatePipeline("h264",[]*webrtc.TrackLocalStaticSample{localTrack}).Start()
}

/*
Use pion mediastream to read in frames
*/
