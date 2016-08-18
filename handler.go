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
	"strconv"
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
		QueryParam      *QueryParam
	}

	QueryParam struct {
		Count  bool
		Sort   string
		Offset int
		Limit  int
		Id     []string
		Field  []string
		Embed  []string
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

	if c.Request().Method() == "GET" {
		h.setQueryParam(c)
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

func (h *Handler) setQueryParam(c echo.Context) {
	qp := new(QueryParam)
	qs := c.QueryParams()

	if param, ok := qs["count"]; ok {
		qp.Count, _ = strconv.ParseBool(param[0])
	}

	if param, ok := qs["embed"]; ok {
		qp.Embed = strings.Split(param[0], ",")
	}

	if param, ok := qs["field"]; ok {
		qp.Field = strings.Split(param[0], ",")
	}

	if param, ok := qs["id"]; ok {
		qp.Id = strings.Split(param[0], ",")
	}

	if param, ok := qs["per_page"]; ok {
		qp.Limit, _ = strconv.Atoi(param[0])
	}

	if param, ok := qs["page"]; ok {
		limit := 10
		page, _ := strconv.Atoi(param[0])

		if qp.Limit != 0 {
			limit = qp.Limit
		}

		qp.Offset = (page - 1) * limit
	}

	if param, ok := qs["sort"]; ok {
		sort := param[0]
		order := "asc"
		if string(sort[0]) == "-" {
			sort = strings.Replace(sort, "-", "", -1)
			order = "desc"
		}

		qp.Sort = fmt.Sprintf("%s %s", sort, order)
	}

	h.QueryParam = qp
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
