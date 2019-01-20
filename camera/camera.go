package camera

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/dchote/gopicamera/config"

	"github.com/mattn/go-mjpeg"
	"gocv.io/x/gocv"
)

var (
	deviceID int
	err      error
	camera   *gocv.VideoCapture

	Stream *mjpeg.Stream
)

const captureWidth = 640
const captureHeight = 480
const captureFPS = 5

func StartCamera() {
	deviceID = config.Config.Camera.DeviceID

	camera, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}

	camera.Set(gocv.VideoCaptureFrameWidth, captureWidth)
	camera.Set(gocv.VideoCaptureFrameHeight, captureHeight)
	camera.Set(gocv.VideoCaptureFPS, captureFPS)

	//camera.Set(gocv.VideoCaptureFormat, 5)
	//defer camera.Close()

	// create the mjpeg stream
	Stream = mjpeg.NewStreamWithInterval(50 * time.Millisecond)

	// start capturing
	go CaptureVideo()
}

func StopCamera() {
	Stream.Close()
	camera.Close()
}

func CaptureVideo() {
	frame := gocv.NewMat()
	defer frame.Close()

	var fistFrame = true

	for {
		if Stream.NWatch() > 0 || fistFrame {
			if ok := camera.Read(&frame); !ok {
				fmt.Printf("Device closed: %v\n", deviceID)
				return
			}

			if frame.Empty() {
				fmt.Printf("Empty image: %v\n", deviceID)
				continue
			}

			// working image
			img := frame.Clone()
			if config.Config.Camera.FlipHorizontal == true && config.Config.Camera.FlipVertical == true {
				gocv.Flip(frame, &img, -1)
				img.CopyTo(&frame)
			} else if config.Config.Camera.FlipHorizontal == true {
				gocv.Flip(frame, &img, 0)
				img.CopyTo(&frame)
			} else if config.Config.Camera.FlipVertical == true {
				gocv.Flip(frame, &img, 1)
				img.CopyTo(&frame)
			}

			var frameWidth = captureWidth
			var frameHeight = captureHeight

			if config.Config.Camera.Rotate == 90 {
				gocv.Rotate(frame, &img, 0)
				img.CopyTo(&frame)

				frameWidth = captureHeight
				frameHeight = captureWidth
			} else if config.Config.Camera.Rotate == 180 {
				gocv.Rotate(frame, &img, 1)
				img.CopyTo(&frame)
			} else if config.Config.Camera.Rotate == 270 {
				gocv.Rotate(frame, &img, 2)
				img.CopyTo(&frame)

				frameWidth = captureHeight
				frameHeight = captureWidth
			}

			if config.Config.Camera.ShowDateTime == true {
				var dateTimePosition image.Point

				if config.Config.Camera.DateTimePosition == "top_right" {
					dateTimePosition = image.Pt(frameWidth-200, 30)
				} else if config.Config.Camera.DateTimePosition == "bottom_right" {
					dateTimePosition = image.Pt(frameWidth-200, frameHeight-20)
				} else if config.Config.Camera.DateTimePosition == "bottom_left" {
					dateTimePosition = image.Pt(20, frameHeight-20)
				} else {
					dateTimePosition = image.Pt(20, 30)
				}

				currentTime := time.Now()
				gocv.PutText(&frame, currentTime.Format("2006-01-02 15:04:05"), dateTimePosition, gocv.FontHersheyPlain, 1, color.RGBA{255, 255, 255, 0}, 2)
			}

			img.Close()

			// encode our processed frame as a JPEG for the MJPEG stream
			buf, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, frame, []int{gocv.IMWriteJpegQuality, 65})
			if err != nil {
				fmt.Printf("error encoding: %v\n", deviceID)
				continue
			}

			Stream.Update(buf)

			fistFrame = false
		}

		// lessen the load a little
		time.Sleep(25 * time.Millisecond)
	}
}
