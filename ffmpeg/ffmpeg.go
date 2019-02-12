package ffmpeg

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"

	//"github.com/dchote/gopicamera/config"

	"github.com/dchote/gmf"
	"github.com/mattn/go-mjpeg"
)

var (
	DeviceName string
	err        error

	CaptureWidth  int
	CaptureHeight int
	CaptureFPS    int

	PixelFormat int32

	Stream *mjpeg.Stream
)

func StartCamera() {
	CaptureWidth = 640
	CaptureHeight = 480
	CaptureFPS = 30

	PixelFormat = gmf.AV_PIX_FMT_YUVJ422P

	// create the mjpeg stream
	Stream = mjpeg.NewStreamWithInterval(15 * time.Millisecond)

	inputCtx := gmf.NewCtx()
	defer inputCtx.CloseInputAndRelease()

	log.Printf("running on %s", runtime.GOOS)

	if runtime.GOOS == "darwin" {
		inputCtx.SetInputFormat("avfoundation")
		DeviceName = "default"
	} else {
		inputCtx.SetInputFormat("video4linux2")
		DeviceName = "/dev/video0"
	}

	err := inputCtx.OpenInputWithOptions(DeviceName, []gmf.Pair{
		{Key: "pixel_format", Val: "uyvy422"},
		{Key: "video_size", Val: fmt.Sprintf("%dx%d", CaptureWidth, CaptureHeight)},
		{Key: "framerate", Val: strconv.Itoa(CaptureFPS)},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	srcVideoStream, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_VIDEO)
	if err != nil {
		log.Println("No video stream found")
		return
	}

	codec, err := gmf.FindEncoder("mjpeg")
	if err != nil {
		log.Fatal(err)
		return
	}

	cc := gmf.NewCodecCtx(codec)
	defer gmf.Release(cc)

	cc.SetPixFmt(PixelFormat)
	cc.SetWidth(CaptureWidth)
	cc.SetHeight(CaptureHeight)
	cc.SetTimeBase(gmf.AVR{1, 1})

	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if err := cc.Open(nil); err != nil {
		log.Fatal(err)
		return
	}

	swsCtx := gmf.NewSwsCtx(srcVideoStream.CodecCtx(), cc, gmf.SWS_FAST_BILINEAR)
	defer gmf.Release(swsCtx)

	dstFrame := gmf.NewFrame().
		SetWidth(CaptureWidth).
		SetHeight(CaptureHeight).
		SetFormat(PixelFormat) // see above
	defer gmf.Release(dstFrame)

	if err := dstFrame.ImgAlloc(); err != nil {
		log.Fatal(err)
		return
	}

	for packet := range inputCtx.GetNewPackets() {
		if Stream == nil {
			break
		}

		if packet.StreamIndex() != srcVideoStream.Index() {
			// skip non video streams
			continue
		}

		ist, err := inputCtx.GetStream(packet.StreamIndex())

	decode:

		frame, err := packet.Frames(ist.CodecCtx())
		if err != nil {
			// Retry if EAGAIN
			if err.Error() == "Resource temporarily unavailable" {
				goto decode
			}
			log.Fatal(err)
		}

		swsCtx.Scale(frame, dstFrame) // TODO I really want to get rid of this!

		p, err := dstFrame.Encode(cc)
		if err != nil {
			gmf.Release(p)
			log.Fatal(err)
			return
		}

		Stream.Update(p.Data())

		gmf.Release(p)
		gmf.Release(frame)
		gmf.Release(packet)
	}
}

func StopCamera() {
	Stream.Close()
	Stream = nil
}
