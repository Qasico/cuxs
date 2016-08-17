package cuxs

import (
	"gopkg.in/go-playground/validator.v8"
	"strings"
	"github.com/fatih/structs"
	"github.com/qasico/cuxs/helper"
)

type(
	Request struct {
		Response *Response
		Validate *validator.Validate
	}
)

func NewRequest(r *Response) *Request {
	req := new(Request)
	req.Validate = validator.New(&validator.Config{TagName: "validate"})
	req.Response = r

	return req
}

func (r *Request) Bind(target interface{}) (fields []string, err error) {
	if err = r.Response.Context.Bind(target); err == nil {
		err = r.Validate.Struct(target)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			r.Response.SetValidation(r.ValidationError(errs))
		}

		fields = r.InputKeys(target)
	}

	return
}

func (r *Request) InputKeys(h interface{}) []string {
	var objmap map[string]string

	rm := structs.Map(h)
	r.Response.Context.Bind(&objmap)
	keys := make([]string, 0, len(objmap))
	for k := range objmap {
		k = helper.CamelCase(k)
		if _, ok := rm[k]; ok {
			keys = append(keys, k)
		}
	}

	return keys
}

func (r *Request) ValidationError(errs validator.ValidationErrors) (v []ResponseValidation) {
	for _, e := range errs {
		x := ResponseValidation{Field: strings.ToLower(e.Field), Message: e.Tag}
		v = append(v, x)
	}

	return
}