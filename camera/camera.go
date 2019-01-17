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
		// lessen the load a little
		time.Sleep(100 * time.Millisecond)

		if ok := camera.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		// TODO logic for control loop feedback

		// write video frame as jpeg to MJPEG stream
		buf, _ := gocv.IMEncode(".jpg", img)
		Stream.UpdateJPEG(buf)
	}
}
