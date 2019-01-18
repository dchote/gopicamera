package camera

import (
	"fmt"
	"time"

	"github.com/dchote/gopicamera/config"

	"github.com/hybridgroup/mjpeg"
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
	camera.Close()
}

func CaptureVideo() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := camera.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// write video frame as jpeg to MJPEG stream
		//gocv.IMEncode(".jpg", img)
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			continue
		}

		Stream.UpdateJPEG(buf)
		// lessen the load a little
		time.Sleep(50 * time.Millisecond)

	}
}
