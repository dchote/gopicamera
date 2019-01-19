package camera

import (
	"fmt"
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

func StartCamera() {
	deviceID = config.Config.Camera.DeviceID

	camera, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	//defer camera.Close()

	// create the mjpeg stream
	Stream = mjpeg.NewStream()

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

	for {
		if Stream.NWatch() > 0 {
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

			if config.Config.Camera.Rotate == 90 {
				gocv.Rotate(frame, &img, 0)
				img.CopyTo(&frame)
			} else if config.Config.Camera.Rotate == 180 {
				gocv.Rotate(frame, &img, 1)
				img.CopyTo(&frame)
			} else if config.Config.Camera.Rotate == 270 {
				gocv.Rotate(frame, &img, 2)
				img.CopyTo(&frame)
			}

			// encode our processed frame as a JPEG for the MJPEG stream
			buf, err := gocv.IMEncode(".jpg", frame)
			if err != nil {
				fmt.Printf("error encoding: %v\n", deviceID)
				continue
			}

			Stream.Update(buf)
		}

		// lessen the load a little
		time.Sleep(250 * time.Millisecond)
	}
}
