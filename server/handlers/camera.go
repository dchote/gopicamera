package handlers

import (
	//"encoding/json"
	//"fmt"
	//"log"
	"net/http"
	//"strings"

	"github.com/dchote/gopicamera/config"

	"github.com/labstack/echo"
)

type CameraObject struct {
	APIURL    string `json:"APIURL"`
	CameraURL string `json:"cameraURL"`
	Name      string `json:"name"`
}

type SurveyObjectList struct {
	Cameras []*CameraObject `json:"cameras"`
}

func CameraList() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO actually build dynamic camera list from mDNS

		cameraList := new(SurveyObjectList)

		camera := &CameraObject{
			APIURL:    config.Config.API.APIURL,
			CameraURL: config.Config.API.CameraURL,
			Name:      config.Config.Camera.Name,
		}

		cameraList.Cameras = append(cameraList.Cameras, camera)

		return c.JSON(http.StatusOK, cameraList)
	}
}
