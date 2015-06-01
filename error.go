package sleepy

type Error interface {
	Error() string
	StatusCode() int
	Message() string
}

type sleepyInternalError struct {
	code    int
	err     error
	message string
}

func (sie *sleepyInternalError) Error() string {
	return sie.Error()
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
