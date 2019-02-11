package ffmpeg

import (
	//"fmt"
	"log"
	"time"

	"github.com/dchote/gopicamera/config"

	"github.com/3d0c/gmf"
	"github.com/mattn/go-mjpeg"
)

var (
	deviceID int
	err      error

	Stream *mjpeg.Stream
)

const captureWidth = 640
const captureHeight = 480
const captureFPS = 5

func StartCamera() {
	deviceID = config.Config.Camera.DeviceID

	// create the mjpeg stream
	Stream = mjpeg.NewStreamWithInterval(50 * time.Millisecond)

	// hardcoding mac stuff for debug
	inputCtx, err := gmf.NewInputCtxWithFormatName("LG UltraFine Display Camera", "avfoundation")
	defer inputCtx.CloseInputAndRelease()

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

	cc.SetPixFmt(gmf.AV_PIX_FMT_YUVJ420P)
	cc.SetWidth(captureWidth)
	cc.SetHeight(captureHeight)
	cc.SetTimeBase(srcVideoStream.CodecCtx().TimeBase().AVR())

	if codec.IsExperimental() {
		cc.SetStrictCompliance(gmf.FF_COMPLIANCE_EXPERIMENTAL)
	}

	if err := cc.Open(nil); err != nil {
		log.Fatal(err)
		return
	}

	swsCtx := gmf.NewSwsCtx(srcVideoStream.CodecCtx(), cc, gmf.SWS_BICUBIC)
	defer gmf.Release(swsCtx)

	dstFrame := gmf.NewFrame().
		SetWidth(captureWidth).
		SetHeight(captureHeight).
		SetFormat(gmf.AV_PIX_FMT_YUVJ420P) // see above
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

		swsCtx.Scale(frame, dstFrame)

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
