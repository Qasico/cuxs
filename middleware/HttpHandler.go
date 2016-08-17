package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

type ResponseFormat struct {
	Code    int           `json:"code,omitempty";xml:"code,omitempty"`
	Message interface{}   `json:"message,omitempty";xml:"message,omitempty"`
}

func HTTPHandler(err error, c echo.Context) {
	var r ResponseFormat

	r.Code = http.StatusInternalServerError
	r.Message = http.StatusText(r.Code)

	if he, ok := err.(*echo.HTTPError); ok {
		r.Code = he.Code
		r.Message = he.Message
	}

	if !c.Response().Committed() {
		if c.Request().Method() == "HEAD" {
			c.NoContent(r.Code)
		} else {
			c.JSON(r.Code, r)
		}
	}
}