package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type ConfigStruct struct {
	Camera struct {
		DeviceID int `json:"deviceID"`
	} `json:"camera"`
	Server struct {
		ListenAddress string `json:"listenAddress"`
		ListenPort    int    `json:"listenPort"`
	} `json:"server"`
	API struct {
		APIURL string `json:"APIURL"`
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

	var apiListenAddress string

	if Config.Server.ListenAddress == "" {
		Config.Server.ListenAddress = GetLocalIP()
	}

	// we MUST provide a valid IP to connect to for API purposes
	if Config.Server.ListenAddress == "0.0.0.0" {
		apiListenAddress = GetLocalIP()
	} else {
		apiListenAddress = Config.Server.ListenAddress
	}

	// build API endpoint URL
	Config.API.APIURL = "http://" + apiListenAddress + ":" + strconv.Itoa(Config.Server.ListenPort) + "/v1"

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
