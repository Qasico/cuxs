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
	"reflect"
	"errors"
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
	h.Validate.RegisterValidation("encrypted", Validencrypted)

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

func Validencrypted(v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value, field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string) bool {

	i, err := strconv.Atoi(field.String())
	if err != nil {
		return false
	}

	val := ((0x0000FFFF & i) << 16) + ((0xFFFF0000 & i) >> 16)
	if val < 1 {
		return false
	}

	return true
}

func (h *Handler) validateRequest(req interface{}) (err error) {
	if err = h.Validate.Struct(req); err != nil {
		h.ValidationError(err.(validator.ValidationErrors))
	} else {
		h.requestKeys(req)
	}

	return
}

func (r *Handler) ValidationError(errs validator.ValidationErrors) {
	for _, e := range errs {
		x := response.ErrorValidation{Field: helper.SnakeCase(e.Field), Message: e.Tag}
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
	qp.Limit = 10
	qp.Offset = 0

	if param, ok := qs["count"]; ok && param[0] != "" {
		qp.Count, _ = strconv.ParseBool(param[0])
	}

	if param, ok := qs["embed"]; ok && param[0] != "" {
		qp.Embed = strings.Split(param[0], ",")
	}

	if param, ok := qs["field"]; ok && param[0] != "" {
		qp.Field = strings.Split(param[0], ",")
	}

	if param, ok := qs["id"]; ok && param[0] != "" {
		qp.Id = strings.Split(param[0], ",")
	}

	if param, ok := qs["per_page"]; ok && param[0] != "" {
		qp.Limit, _ = strconv.Atoi(param[0])
	}

	if param, ok := qs["page"]; ok && param[0] != "" {
		limit := 10
		page, _ := strconv.Atoi(param[0])

		if qp.Limit != 0 {
			limit = qp.Limit
		}

		qp.Offset = (page - 1) * limit
	}

	if param, ok := qs["sort"]; ok && param[0] != "" {
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
			if h.Response.Code > 300 {
				h.Response.SetCode(http.StatusOK)
			}

			h.Response.Status = response.StatusSuccess
			h.Response.Message = nil

			if h.ResponseHandler != nil {
				h.Response.Data = h.ResponseHandler
			}

			h.FilterResponse()
		}
	}

	return h.Context.JSON(h.Response.Code, h.Response)
}

func (h *Handler) FilterResponse() {
	// run only GET requests
	if h.Context.Request().Method() == "GET" && len(h.QueryParam.Field) > 0 && h.Response.Data != nil && h.Response.Status == response.StatusSuccess {
		// filter the Response.Data
		d := h.Response.Data
		switch reflect.TypeOf(d).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(d)
			var result []interface{}
			for i := 0; i < s.Len(); i++ {
				m := make(map[string]interface{})
				x := s.Index(i)
				for _, fname := range h.QueryParam.Field {
					if fname == "id" {
						fname = "id_e"
					}

					m[fname] = x.FieldByName(helper.CamelCase(fname)).Interface()
				}
				result = append(result, m)
			}
			h.Response.Data = result
		case reflect.Ptr:
			rm := structs.Map(d)
			x := make(map[string]interface{})
			for _, fname := range h.QueryParam.Field {
				if fname == "id" {
					fname = "id_e"
				}
				kk := helper.CamelCase(fname)
				if _, ok := rm[kk]; ok {
					x[fname] = rm[kk]
				}
			}

			h.Response.Data = x
		}
	}
}

func (h *Handler) SetCreated(d interface{}) {
	h.Response.SetCode(response.StatusCreated)
	h.Response.SetData(d)
}

func (h *Handler) Valid(name string, field interface{}, rule string) error {
	if err := h.Validate.Field(field, rule); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			x := response.ErrorValidation{Field: helper.SnakeCase(name), Message: e.Tag}
			h.Response.Errors = append(h.Response.Errors, x)
		}

		return errors.New("Validation Failed")
	}

	return nil
}

func (q *QueryParam) IsEmbed(field string) bool {
	if q != nil && len(q.Embed) > 0 {
		for _, a := range q.Embed {
			if a == field {
				return true
			}
		}
	}

	return false
}

func SetResponseMessage(err string) {
	ApiHandler.Response.SetCode(response.StatusBadRequest)
	ApiHandler.Response.Message = err
}

func SetHandler(res ResponseHandler) {
	ApiHandler.ResponseHandler = &res
}

func SetErrorValidate(field string, message string) {
	x := response.ErrorValidation{Field: helper.SnakeCase(field), Message: message}

	ApiHandler.Response.Errors = append(ApiHandler.Response.Errors, x)
}
