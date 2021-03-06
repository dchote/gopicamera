package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	//"strconv"
)

type ConfigStruct struct {
	Camera struct {
		DeviceID         int    `json:"deviceID"`
		Name             string `json:"name"`
		Rotate           int    `json:"rotate"`
		FlipHorizontal   bool   `json:"flipHorizontal"`
		FlipVertical     bool   `json:"flipVertical"`
		ShowDateTime     bool   `json:"showDateTime"`
		DateTimePosition string `json:"dateTimePosition"`
	} `json:"camera"`
	Server struct {
		ListenAddress string `json:"listenAddress"`
		ListenPort    int    `json:"listenPort"`
	} `json:"server"`
	API struct {
		APIURL    string `json:"APIURL"`
		CameraURL string `json:"cameraURL"`
	}
}

var (
	ConfigFile string
	Config     *ConfigStruct
)

func LoadConfig(file string) (*ConfigStruct, error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Config)

	if err != nil {
		return nil, err
	}

	if Config.Server.ListenAddress == "" {
		Config.Server.ListenAddress = GetLocalIP()
	}

	// build API endpoint URL
	Config.API.APIURL = "/v1"
	Config.API.CameraURL = "/camera.mjpeg"

	return Config, err
}

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
