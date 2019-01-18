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
	img := gocv.NewMat()
	defer img.Close()

	for {
		if Stream.NWatch() > 0 {
			if ok := camera.Read(&img); !ok {
				fmt.Printf("Device closed: %v\n", deviceID)
				return
			}

			if img.Empty() {
				fmt.Printf("Empty image: %v\n", deviceID)
				continue
			}

			// write video frame as jpeg to MJPEG stream
			//gocv.IMEncode(".jpg", img)
			buf, err := gocv.IMEncode(".jpg", img)
			if err != nil {
				fmt.Printf("error encoding: %v\n", deviceID)
				continue
			}

			Stream.Update(buf)
		}

		// lessen the load a little
		time.Sleep(50 * time.Millisecond)
	}
}
