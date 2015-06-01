package main

type Error struct {
	responseCode int
	message      string
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) StatusCode() int {
	return e.responseCode
}

var (
	ErrInternal      = &Error{responseCode: 500, message: "Server error."}
	ErrLogin         = &Error{responseCode: 401, message: "Please login."}
	ErrAuthorization = &Error{responseCode: 403, message: "Unauthorized request."}
)

func NewRequestError(msg string) *Error {
	return &Error{responseCode: 400, message: msg}
}
