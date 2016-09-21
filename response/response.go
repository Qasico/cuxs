package response

type (
	Attribute struct {
		Code    int               `json:"-"`
		Status  string            `json:"status,omitempty"`
		Message interface{}       `json:"message,omitempty"`
		Data    interface{}       `json:"data,omitempty"`
		Total   int64             `json:"total,omitempty"`
		Errors  []ErrorValidation `json:"errors,omitempty"`
	}

	ErrorValidation struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func (r *Attribute) SetCode(s int) *Attribute {
	r.Code = s

	return r
}

func (r *Attribute) SetData(d interface{}) *Attribute {
	r.Data = d

	return r
}

func (r *Attribute) SetMessage(m interface{}) *Attribute {
	r.Message = m

	return r
}

func (r *Attribute) SetError(field string, message string) *Attribute {
	ev := ErrorValidation{Field: field, Message: message}

	r.Errors = append(r.Errors, ev)
	return r
}
