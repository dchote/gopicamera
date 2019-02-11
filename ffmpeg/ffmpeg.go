package ffmpeg

import (
	//"fmt"
	"log"
	"time"

	"github.com/dchote/gopicamera/config"

	"github.com/mattn/go-mjpeg"

	"github.com/3d0c/gmf"
)

var (
	deviceID int
	err      error

	ffmpegContext  *gmf.FmtCtx
	srcVideoStream *gmf.Stream

	Stream *mjpeg.Stream
)

const captureWidth = 640
const captureHeight = 480
const captureFPS = 5

func StartCamera() {
	deviceID = config.Config.Camera.DeviceID

	// hardcoding mac stuff for debug
	ffmpegContext, err := gmf.NewInputCtxWithFormatName("LG UltraFine Display Camera", "avfoundation")

	if err != nil {
		log.Fatal(err)
	}

	srcVideoStream, err := ffmpegContext.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	log.Printf("stream: %s", srcVideoStream)
	if err != nil {
		log.Println("No video stream found")
	}

	codec, err := gmf.FindEncoder("mjpeg")
	log.Printf("codec: %s", codec)
	if err != nil {
		log.Fatal(err)
	}

	// create the mjpeg stream
	Stream = mjpeg.NewStreamWithInterval(50 * time.Millisecond)

	// start capturing
	go CaptureVideo()
}

func StopCamera() {
	Stream.Close()
	ffmpegContext.CloseInputAndRelease()
}

func CaptureVideo() {
	var fistFrame = true

	for {
		if Stream.NWatch() > 0 || fistFrame {

			//Stream.Update(buf)

			fistFrame = false
		}

		// lessen the load a little
		time.Sleep(25 * time.Millisecond)
	}
}
