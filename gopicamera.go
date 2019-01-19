package main

import (
	//"flag"
	//"fmt"
	//"io/ioutil"
	"log"
	"os"
	"os/signal"
	//"strings"
	"syscall"
	"time"

	"github.com/dchote/gopicamera/camera"
	"github.com/dchote/gopicamera/config"
	"github.com/dchote/gopicamera/server"

	"github.com/GeertJohan/go.rice"
	"github.com/docopt/docopt-go"
)

const VERSION = "0.0.1"

var (
	err          error
	staticAssets *rice.Box
)

func cliArguments() {
	usage := `
Usage: gopicamera [options]

Options:
  -c, --config=<json>           Specify config file [default: ./config.json]
	-d, --camera-device=<device>  Specify the devide id of the camera [default: 0]
	-a, --assets=<dir>            Static assets directory [default: ./frontend/dist]
  -h, --help                    Show this screen.
  -v, --version                 Show version.
`
	args, _ := docopt.ParseArgs(usage, os.Args[1:], VERSION)

	config.ConfigFile, _ = args.String("--config")

	_, err = config.LoadConfig(config.ConfigFile)
	if err != nil {
		log.Fatalf("Unable to load "+config.ConfigFile+" ERROR=", err)
	}

	// override the camera device if specified
	cameraDeviceOverride, _ := args.Int("--camera-device")
	if cameraDeviceOverride > 0 {
		config.Config.Camera.DeviceID = cameraDeviceOverride
	}

	staticAssetsDirectory, _ := args.String("--assets")
	if staticAssetsDirectory != config.Config.Server.StaticAssetsDirectory && staticAssetsDirectory != "frontend/dist" {
		config.Config.Server.StaticAssetsDirectory = staticAssetsDirectory
	}

	log.Printf("Config: %+v", config.Config)
}

func main() {
	cliArguments()

	log.Printf("Loading static assets from %s", config.Config.Server.StaticAssetsDirectory)
	staticAssets, err = rice.FindBox(config.Config.Server.StaticAssetsDirectory)
	if err != nil {
		log.Fatalf("Static assets not found. Build them with npm first.")
	}

	// ignore sigpipe that happens when the mjpeg stream is terminated by the client
	signal.Ignore(syscall.SIGPIPE)

	// start camera
	camera.StartCamera()

	// start the webserver
	go server.StartServer(*config.Config, staticAssets)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("Shutting down")

	// shut down listener, with a hard timeout
	server.StopServer()
	camera.StopCamera()

	// extra grace time
	time.Sleep(time.Second)

	os.Exit(0)
}
