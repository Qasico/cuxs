package response

const (
	StatusOK = 200
	StatusCreated = 201
	StatusBadRequest = 400
	StatusUnauthorized = 401
	StatusNotFound = 404
	StatusUnprocessableEntry = 422
	StatusInternalServerError = 500
	StatusFailed = "fail"
	StatusSuccess = "success"
)

var statusText = map[int]string{
	StatusOK:                           "OK",
	StatusCreated:                      "Created",
	StatusBadRequest:                   "Bad Request",
	StatusUnauthorized:                 "Unauthorized",
	StatusNotFound:                     "Not Found",
	StatusUnprocessableEntry:           "Validation Failed",
	StatusInternalServerError:          "Internal Server Error",
}

// StatusText returns a text for the HTTP status code. It returns the empty
// string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
