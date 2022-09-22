package service

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"

	"github.com/pion/webrtc/v3"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Stream struct {
	Cmd      *exec.Cmd
	Listener *net.UDPConn
}

func NewStream() *Stream {
	return &Stream{}
}

// Plays a stream
func (s *Stream) Play(movieFile string, track *webrtc.TrackLocalStaticRTP) error {
	//start udp listeners
	adpAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 5004}
	udpconn, err := net.ListenUDP("udp", &adpAddr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer func() {
		errC := udpconn.Close()
		if errC != nil {
			panic(errC)
		}
	}()
	//start ffmpeg streamers

	//Start ffmpeg process to output to RTP

	buf := bytes.NewBuffer(nil)
	inputVid := ffmpeg.Input(movieFile).Video()
	/* inputAud := ffmpeg.Input(dir+filename).Audio().Output("rtp://127.0.0.1:5006?pkt_size=1200",
	ffmpeg.KwArgs{"acodec": "copy", "f": "rtp"}) */
	cmd := inputVid.
		Output("rtp://127.0.0.1:5004?pkt_size=1200",
			ffmpeg.KwArgs{"c:v": "libx264", "f": "rtp", "g": "10", "tune": "zerolatency", "r": "24", "pix_fmt": "yuv420p", "filter:v": "setpts=2.0*PTS"}).
		WithOutput(buf, os.Stdout).
		Compile()

	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting stream.")
		fmt.Println(err)
	}

	//store cmd.exe controller in Stream object
	s.Cmd = cmd
	//store udp listener in Streamobject
	s.Listener = udpconn
	//Start reading from udp port and writing to
	//stream from RTP connection to webrtc
	inRTPpack := make([]byte, 1600)
	for {
		n, _, errRead := s.Listener.ReadFrom(inRTPpack)
		if errRead != nil {
			fmt.Println("Error in reading RTP Packets: ")
			panic(errRead)
		}

		_, errWrite := track.Write(inRTPpack[:n])
		if errWrite == io.ErrClosedPipe {
			fmt.Println(errWrite)
			break
		}
	}
	//webrtc track
	return nil
}

// Stops a stream
func (s *Stream) Stop() error {

	//using *exec.cmd stop stream manually
	err := s.Cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error in killing process")
		panic(err)
	}
	return nil
}

/*
	On stream request, a stream object is created
	containing the compiled cmd.exe object which controls
	the process.
	When the stream object is created,its stored in the agent object
	which made the request.
	The stream is then started

	Then user can also stop stream via an endpoint to which the agent
	object can also control where a cmd.Kill process is called on the
	playingstream
*/
