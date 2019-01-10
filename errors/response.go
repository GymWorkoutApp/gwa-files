package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Response error response
type Response struct {
	URI         string      `json:"-"`
	StatusCode  int         `json:"-"`
	Header      http.Header `json:"-"`
	Message  	interface{} `json:"message"`
	Internal 	error       `json:"-"`
}

// NewResponse create the response pointer
func NewResponse(err error, statusCode int, message string) *Response {
	return &Response{
		Internal:      err,
		StatusCode: statusCode,
		Message: message,
	}
}

// NewResponse create the response pointer
func NewResponseByError(err error) *Response {
	return &Response{
		Internal:      err,
		StatusCode: StatusCodes[err],
		Message: Descriptions[err],
	}
}

// SetHeader sets the header entries associated with key to
// the single element value.
func (r *Response) SetHeader(key, value string) {
	if r.Header == nil {
		r.Header = make(http.Header)
	}
	r.Header.Set(key, value)
}

func (he *Response) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.StatusCode, he.Message)
}

func (he *Response) SetInternal(err error) *Response {
	he.Internal = err
	return he
}

// https://tools.ietf.org/html/rfc6749#section-5.2
var (
	ErrFileNotFound = 			errors.New("file_not_found")
	ErrJsonMarshal = 			errors.New("json_marshal_error")
	ErrYouAreNotAuthenticated = errors.New("not_authenticated")

)

// Descriptions error description
var Descriptions = map[error]string{
	ErrFileNotFound: 			"The file not found.",
	ErrJsonMarshal:  			"The object is not serializable to json.",
	ErrYouAreNotAuthenticated:  "You are not authorized for this operation.",
}

// StatusCodes response error HTTP status code
var StatusCodes = map[error]int{
	ErrFileNotFound: 			400,
	ErrJsonMarshal:  			422,
	ErrYouAreNotAuthenticated:  401,
}
