package sleepy

type Error struct {
	HttpCode int    `json:"-"`
	Err      string `json:"error"`
	Msg      string `json:"message"`
	Code     int    `json:"code"`
}

func ErrInternal(err string) *Error {
	return &Error{500, err, "", ERR_INTERNAL}
}

func ErrBadRequest(err string, msg string, code int) *Error {
	return &Error{422, err, msg, code}
}

const (
	ERR_INTERNAL = 1000 + iota
	ERR_PARSE_REQUEST
	ERR_FIELD_MISSING
	ERR_MOD_RO_FIELD
)
