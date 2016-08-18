package cuxs

import (
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v8"
	"github.com/qasico/cuxs/response"
	"github.com/fatih/structs"
	"github.com/qasico/cuxs/helper"
	"strings"
	"net/http"
	"fmt"
)

type(
	RequestHandler interface{}

	ResponseHandler interface{}

	Handler struct {
		Context         echo.Context
		Response        *response.Attribute
		Validate        *validator.Validate
		RequestHandler  *RequestHandler
		ResponseHandler *ResponseHandler
		RequestInput    []string
	}
)

var ApiHandler *Handler

func (h *Handler) Prepare(c echo.Context, req RequestHandler) (hr *Handler, err error) {
	h.Validate = validator.New(&validator.Config{TagName: "validate"})
	h.Response = &response.Attribute{Code: response.StatusBadRequest, Status: response.StatusFailed, Message: response.StatusText(response.StatusBadRequest)}
	h.Context = c
	h.ResponseHandler = nil

	if req != nil {
		h.Context.Bind(&req)
		if err = h.validateRequest(req); err == nil {
			h.RequestHandler = &req
		}
	}

	ApiHandler = h

	return h, err
}

func (h *Handler) validateRequest(req interface{}) (err error) {
	if err = h.Validate.Struct(req); err != nil {
		h.ValidationError(err.(validator.ValidationErrors))
	} else {
		h.requestKeys(req)
		fmt.Println("OK")
	}

	return
}

func (r *Handler) ValidationError(errs validator.ValidationErrors) {
	for _, e := range errs {
		x := response.ErrorValidation{Field: strings.ToLower(e.Field), Message: e.Tag}
		r.Response.Errors = append(r.Response.Errors, x)
	}
}

func (h *Handler) requestKeys(i interface{}) {
	var objmap map[string]string

	rm := structs.Map(i)
	h.Context.Bind(&objmap)
	keys := make([]string, 0, len(objmap))
	for k := range objmap {
		kk := helper.CamelCase(k)
		if _, ok := rm[kk]; ok {
			keys = append(keys, k)
		}
	}

	h.RequestInput = keys
}

func (h *Handler) Serve(err error) error {
	// check if errors has contain data
	if len(h.Response.Errors) > 0 {
		h.Response.SetCode(response.StatusUnprocessableEntry)
		h.Response.SetMessage(response.StatusText(response.StatusUnprocessableEntry))
	} else {
		if err != nil {
			h.Response.SetCode(response.StatusBadRequest)
			h.Response.SetMessage(err.Error())
			h.Response.Data = nil
		} else {
			h.Response.SetCode(http.StatusOK)
			h.Response.Status = response.StatusSuccess
			h.Response.Message = nil
			if h.ResponseHandler != nil {
				h.Response.SetData(h.ResponseHandler, 0)
			}
		}
	}

	return h.Context.JSON(h.Response.Code, h.Response)
}

func SetHandler(res ResponseHandler) {
	ApiHandler.ResponseHandler = &res
}

func SetErrorValidate(field string, message string) {
	x := response.ErrorValidation{Field: strings.ToLower(field), Message: message}

	ApiHandler.Response.Errors = append(ApiHandler.Response.Errors, x)
}

func GetInputKeys() []string {
	return ApiHandler.RequestInput
}
