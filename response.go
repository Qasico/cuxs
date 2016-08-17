package cuxs

import (
	"net/http"

	"github.com/labstack/echo"
)

type(
	Response struct {
		Context echo.Context
		Format  *ResponseFormat
	}

	ResponseFormat struct {
		Code       int                        `json:"code,omitempty";xml:"code,omitempty"`
		Data       interface{}                `json:"data,omitempty";xml:"data,omitempty"`
		Total      int64                      `json:"total,omitempty";xml:"total,omitempty"`
		Validation []ResponseValidation       `json:"validation,omitempty";xml:"validation,omitempty"`
		Message    interface{}                `json:"message,omitempty";xml:"message,omitempty"`
	}

	ResponseValidation struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func NewResponse(c echo.Context) *Response {
	r := new(Response)

	r.Context = c
	r.Format = new(ResponseFormat)

	// default is bad request
	r.SetCode(http.StatusBadRequest)

	return r
}

func (r *Response) SetCode(s int) {
	r.Format.Code = s
	r.Format.Message = http.StatusText(s)
}

func (r *Response) SetData(d interface{}) {
	r.SetCode(http.StatusOK)
	r.Format.Data = d
}

func (r *Response) SetTotal(t int64) {
	r.Format.Total = t
}

func (r *Response) SetMessage(s interface{}) {
	r.Format.Message = s
}

func (r *Response) SetValidation(s []ResponseValidation) {
	r.Format.Validation = s
	r.SetCode(http.StatusBadRequest)
}

// check if response contains error validations
func (r *Response) isValid() bool {
	return (len(r.Format.Validation) == 0)
}

func (r *Response) Response(code int, d interface{}, msg string) error {

	r.SetCode(code)
	r.SetData(d)
	r.SetMessage(msg)

	return r.Show(r.Format)
}

func (r *Response) Success(data interface{}, total int64) error {
	r.SetCode(http.StatusOK)
	r.SetData(data)

	if total > 0 {
		r.SetTotal(total)
	}

	return r.Show(r.Format)
}

func (r *Response) Error(e error) error {
	r.SetCode(http.StatusBadRequest)
	r.SetMessage(e.Error())
	r.Format.Data = nil

	return r.Show(r.Format)
}

func (r *Response) Serve(err error) error {
	if len(r.Format.Validation) < 1 {
		if err != nil {
			r.SetCode(http.StatusBadRequest)
			r.SetMessage(err.Error())
			r.Format.Data = nil
		} else {
			r.SetCode(http.StatusOK)
		}
	}

	return r.Show(r.Format)
}

func (r *Response) Show(p interface{}) error {
	switch Config.ResponseType {
	case "xml":
		return r.Context.XML(r.Format.Code, p)
	default:
		return r.Context.JSON(r.Format.Code, p)
	}
}

func (r *Response) Write(b []byte) {
	_, e := r.Context.Response().Write(b)
	if e != nil {
		print(e.Error())
	}
}
