package service

import (
	//"bytes"
	//"bytes"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	ffmpeg "github.com/u2takey/ffmpeg-go"

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
	RTPTrack *webrtc.TrackLocalStaticRTP
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
				URLs: []string{"stun:stun3.1.google.com:19302", "stun:stun4.1.google.com:19302"},
			},
		},
	}
	peerconnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return &Agent{Pconnect: peerconnection}, nil
}

func (agent *Agent) SetTrack(typ string) error {
	if typ == "sample" {
		fmt.Println("Local Static Track")
		localtrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion1")
		if err != nil {
			return err
		}

		//add track to peerconnection
		rtpSender, err := agent.Pconnect.AddTrack(localtrack)
		if err != nil {
			panic(err)
		}
		agent.Track = localtrack

		go func() {
			rtcpbuff := make([]byte, 1600)
			for {
				if _, _, errRtcp := rtpSender.Read(rtcpbuff); errRtcp != nil {
					return
				}
			}
		}()
		return nil
	}
	localtrack, errT := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeH264}, "video", "pion2")
	if errT != nil {
		panic(errT)
	}
	rtpSender, errRTP := agent.Pconnect.AddTrack(localtrack)
	if errRTP != nil {
		panic(errRTP)
	}
	agent.RTPTrack = localtrack
	//continuosly read RTCP packets for NACK
	go func() {
		rtcpbuff := make([]byte, 1600)
		for {
			if _, _, errRtcp := rtpSender.Read(rtcpbuff); errRtcp != nil {
				return
			}
		}
	}()
	return nil
}

func (agent *Agent) InitProcess() error {
	//this creates local track
	//for static sample use "sample"
	agent.SetTrack("sample")

	//create onice change listener
	agent.Pconnect.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("\\n New Connection state %v", connectionState)
	})

	agent.Pconnect.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		fmt.Printf("Ice Candidate %v \n", candidate)
		if candidate != nil {
			fmt.Printf("Ice candidate %v", candidate)

			errC := agent.Ws.WriteJSON(candidate)
			if errC != nil {
				panic(errC)
			}
			//agent.Ws.WriteJSON(agent.Pconnect.CurrentLocalDescription())
			//go agent.Ws.ReadJSON(agent.Pconnect.RemoteDescription())

		}
		if candidate == nil {
			fmt.Printf("Null candidate: /n %v", candidate)
			//agent.Ws.WriteJSON(agent.Pconnect.CurrentLocalDescription())
		}
	})

	agent.Pconnect.OnNegotiationNeeded(func() {
		fmt.Println("Negotiation Needed")
	})

	//return track
	return nil
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
		//fmt.Println(buf)
		errtwr := agent.Track.WriteSample(media.Sample{Data: buf[:n]})
		if errtwr != nil {
			fmt.Println(err)
			panic(err)
		}
	}

}

// Use ffmpeg to stream to rtp
// and read from rtp wbertc
func (agent *Agent) StreamRTP() {
	//RTP connection should only be available for
	//the period of streaming, so we close it as
	//soon as streaming is done

	//Opening files
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

	//open RTP connection
	adpAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004}
	udplistener, errUdp := net.ListenUDP("udp", &adpAddr)
	if errUdp != nil {
		panic(errUdp)
	}
	defer func() {
		errC := udplistener.Close()
		if errC != nil {
			panic(errC)
		}
	}()

	//Start ffmpeg process to output to RTP
	go func() {
		buf := bytes.NewBuffer(nil)
		inputVid := ffmpeg.Input(dir + filename).Video()
		/* inputAud := ffmpeg.Input(dir+filename).Audio().Output("rtp://127.0.0.1:5006?pkt_size=1200",
		ffmpeg.KwArgs{"acodec": "copy", "f": "rtp"}) */
		errFF := inputVid.
			Output("rtp://127.0.0.1:5004?pkt_size=1200",
				ffmpeg.KwArgs{"c:v": "libx264", "preset": "ultrafast", "f": "rtp", "error-resilient": "1", "tune": "zerolatency", "r": "30", "pix_fmt": "yuv420p"}).
			WithOutput(buf, os.Stdout).
			Run()
		if errFF != nil {
			panic(errFF)
		}
	}()

	//stream from RTP connection to webrtc
	inRTPpack := make([]byte, 1600)
	for {
		n, _, errRead := udplistener.ReadFrom(inRTPpack)
		if errRead != nil {
			fmt.Println("Error in reading RTP Packets: \\n")
			panic(errRead)
		}

		_, errWrite := agent.RTPTrack.Write(inRTPpack[:n])
		if errWrite == io.ErrClosedPipe {
			return
		}
	}

}
