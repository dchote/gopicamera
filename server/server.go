package server

import (
	"context"
	//"errors"
	"log"
	"net/http"
	//"strings"
	"strconv"
	//"sync"
	"time"

	"github.com/dchote/gopicamera/camera"
	"github.com/dchote/gopicamera/config"
	"github.com/dchote/gopicamera/server/handlers"

	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	e *echo.Echo
)

func StartServer(cfg config.ConfigStruct, assets *rice.Box) {
	if e != nil {
		return
	}

	// instantiate echo instance
	e = echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = 1 * time.Second
	e.Server.WriteTimeout = 0

	// setup middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// prevent caching by client (e.g. Safari)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return next(c)
		}
	})

	if assets != nil {
		assetHandler := http.FileServer(assets.HTTPBox())
		e.GET("/", echo.WrapHandler(assetHandler))
		e.GET("/favicon.ico", echo.WrapHandler(assetHandler))
		e.GET("/css/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/fonts/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
		e.GET("/img/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	}

	// setup API routes
	e.GET("/health", handlers.Health())
	e.GET("/config", handlers.Config())

	// versioned API logic
	e.GET("/v1/camera/list", handlers.CameraList())

	// camera MJPEG stream
	e.GET("/camera.mjpeg", echo.WrapHandler(camera.Stream))

	log.Println("starting server on http://" + cfg.Server.ListenAddress + ":" + strconv.Itoa(cfg.Server.ListenPort))
	e.Logger.Fatal(e.Start(cfg.Server.ListenAddress + ":" + strconv.Itoa(cfg.Server.ListenPort)))
}

func StopServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
