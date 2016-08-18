package middleware

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/qasico/cuxs/response"
)

func HTTPHandler(err error, c echo.Context) {
	var r response.Attribute

	code := response.StatusInternalServerError

	r.Status = response.StatusFailed
	r.Message = http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		r.Message = he.Message
	}

	if !c.Response().Committed() {
		if c.Request().Method() == "HEAD" {
			c.NoContent(code)
		} else {
			c.JSON(code, r)
		}
	}
}