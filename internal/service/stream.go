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
	Playing  bool
}

func NewStream() *Stream {
	return &Stream{Playing: false}
}

// Start udp connection
func (s *Stream) StartUdp(port int) error {
	//start udp listeners
	adpAddr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port}
	udpconn, err := net.ListenUDP("udp", &adpAddr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	s.Listener = udpconn
	return nil
}

// Close Udp connection
func (s *Stream) CloseUdp() error {
	err := s.Listener.Close()
	if err != nil {
		return err
	}

	return nil
}

// Plays a stream
func (s *Stream) Play(movieFile string, track *webrtc.TrackLocalStaticRTP) error {
	//Start ffmpeg process to output to RTP

	buf := bytes.NewBuffer(nil)
	inputVid := ffmpeg.Input(movieFile).Video()
	cmd := inputVid.
		Output("rtp://127.0.0.1:5004?pkt_size=1200",
			ffmpeg.KwArgs{"c:v": "libx264", "f": "rtp", "g": "10", "tune": "zerolatency", "r": "24", "pix_fmt": "yuv420p", "filter:v": "setpts=2.0*PTS"}).
		WithOutput(buf, os.Stdout).
		Compile()

	err := cmd.Start()

	if err != nil {
		fmt.Println("Error starting stream.")
		fmt.Println(err)
		return err
	}

	//store cmd.exe controller in Stream object
	s.Cmd = cmd
	//store udp listener in Streamobject

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

	return nil
}

// Stops a stream
func (s *Stream) Stop() error {
	s.Playing = false
	//using *exec.cmd stop stream manually
	err := s.Cmd.Process.Kill()
	if err != nil {
		fmt.Println("Error in killing process")
		return err
	}
	return nil
}
