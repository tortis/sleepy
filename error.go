package sleepy

type Error interface {
	Error() string
	StatusCode() int
	Message() string
}

type sleepyRequestError struct {
	code    int
	err     error
	message string
}

func (sre *sleepyRequestError) Error() string {
	return sre.err.Error()
}

func (sre *sleepyRequestError) StatusCode() int {
	return sre.code
}

func (sre *sleepyRequestError) Message() string {
	return sre.message
}

func newRequestError(msg string, err error) Error {
	return &sleepyInternalError{
		code:    400,
		err:     err,
		message: msg,
	}
}

type sleepyInternalError struct {
	code    int
	err     error
	message string
}

func (sie *sleepyInternalError) Error() string {
	return sie.err.Error()
}

func (sie *sleepyInternalError) StatusCode() int {
	return sie.code
}

func (sie *sleepyInternalError) Message() string {
	return sie.message
}

func newInternalError(err error) Error {
	return &sleepyInternalError{
		code:    500,
		err:     err,
		message: "Server Error",
	}
}
