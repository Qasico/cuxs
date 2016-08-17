package middleware

import (
	"fmt"
	"time"
	"strconv"

	"github.com/labstack/echo"
	"github.com/qasico/cuxs/log"
)

// Logger returns a middleware that logs HTTP requests.
func HttpLogger() echo.MiddlewareFunc {
	return func(n echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			return logRequest(n, c)
		}
	}
}

func logRequest(hand echo.HandlerFunc, c echo.Context) (err error) {
	start := time.Now()
	msg := "OK"

	if err = hand(c); err != nil {
		msg = err.Error()
		c.Error(err)
	}

	stop := time.Now()
	req := c.Request()
	res := c.Response()
	l := stop.Sub(start).String()

	log.Infof("%3s | %10v | %-8v %-50s %v", getCode(res.Status()), l, getMethod(req.Method()), req.URL().Path(), msg)

	return
}

func getCode(code int) string {
	switch {
	case code >= 200 && code < 300:

		return log.Color.GreenBg(" " + strconv.Itoa(code) + " ", "1")
	case code >= 300 && code < 400:
		return log.Color.YellowBg(" " + strconv.Itoa(code) + " ", "1")
	default:
		return log.Color.RedBg(" " + strconv.Itoa(code) + " ", "1")
	}
}

func getMethod(method string) (m string) {
	switch {
	case method == "GET":
		m = log.Color.Green(fmt.Sprintf("%-8v", method), "1")
	case method == "POST":
		m = log.Color.Blue(fmt.Sprintf("%-8v", method), "1")
	case method == "PUT":
		m = log.Color.Yellow(fmt.Sprintf("%-8v", method), "1")
	case method == "DELETE":
		m = log.Color.Red(fmt.Sprintf("%-8v", method), "1")
	default:
		m = log.Color.Grey(fmt.Sprintf("%-8v", method), "1")
	}

	return
}