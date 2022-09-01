package main

import (
	"fmt"
	"mp4stream/internal/server"

	"github.com/gorilla/mux"
	//"log"
	//"os"
	//ffmpeg "github.com/u2takey/ffmpeg-go"
	//"log"
	//"github.com/pion/webrtc/v3"
)

func main() {
	fmt.Println("FFMPEG")

	//create router instance
	router := mux.NewRouter()
	server := server.NewServer(router)

	server.Listen()
	/* filename := "./assets/DC League of Super-Pets (2022).mp4"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println(file) */
	/* filename := "./assets/DC League of Super-Pets (2022).mp4"
	err := ffmpeg.Input(filename).Output("./assets/out.mp4", ffmpeg.KwArgs{"c:v": "libx265"}).Run()
	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	if err == nil {
		fmt.Println("No error")
	} */
	//fmt.Printf("%v", file)
	/* fmt.Println("Hello Stream")

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.1.google.com:19302", "stun:stun.2.google.com:19302"},
			},
		},
	}

	//creating media engine, will stream .mp4 hence uses h264 codec
	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.RTPCodecParameters{}, webrtc.NewRTPCodecType("h246"))

	//api := webrtc.NewAPI(webrtc.WithMediaEngine(&m))
	//create peer connection
	peerconnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	//create webrtc signal description
	localTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/h264"}, "video", "pion1")
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
	fmt.Printf("Offer: %v", offer) */
	//create answer and establish local SDP and send response to
	//remote peer
	//gst.CreatePipeline("h264",[]*webrtc.TrackLocalStaticSample{localTrack}).Start()
}

/*
Use pion mediastream to read in frames
*/

/*
	Convert video to h264
*/
