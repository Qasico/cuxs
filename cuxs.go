package cuxs

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/qasico/cuxs/log"
	"github.com/qasico/cuxs/middleware"

	mw "github.com/labstack/echo/middleware"
)

var (
	Echo        *echo.Echo
	ContentType string = "JSON"
)

func NewEcho() *echo.Echo {
	Echo = echo.New()
	return Echo
}

func JwtKey() []byte {
	return []byte(Config.JwtHash)
}

func Run() {
	Echo.SetLogger(log.New("-"))
	Echo.SetHTTPErrorHandler(middleware.HTTPHandler)
	Echo.Use(mw.Recover())

	if Config.Runmode == "dev" {
		Echo.Use(middleware.HttpLogger())
		Echo.SetDebug(true)

		// List all declared routes on consoles
		listRoutes()
	}

	log.Infof("Server running on %s", Config.ServerConfig.HTTPAddr)
	Echo.Run(fasthttp.WithTLS(Config.ServerConfig.HTTPAddr, Config.ServerConfig.HTTPSCertFile, Config.ServerConfig.HTTPSKeyFile))
}

func listRoutes() {
	log.Infof("------------------------------------------------------------------------------")
	log.Infof("%-10s | %-50s | %-100s", "METHOD", "URL PATH", "REQ. HANDLER")
	log.Infof("------------------------------------------------------------------------------")
	for _, v := range Echo.Routes() {
		if v.Path[len(v.Path)-1:] != "*" {
			log.Infof("%-10s | %-50s | %-100s", v.Method, v.Path, v.Handler)
		}
	}
	log.Infof("------------------------------------------------------------------------------")
}
