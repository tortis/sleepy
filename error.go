package sleepy

////////////////////////////////////////////////////////////////////////////////
// The error type used by the sleepy library. In the event of an error, this  //
// struct will be used to log the error and write a response to the client.   //
//                                                                            //
// Unfortunately this restricts the errors of sleepy users, as their API      //
// error system must conform to sleepy's error type.                          //
//                                                                            //
// Sleepy also defines some API error codes that are expected to be used by   //
// the library user in addition to their own custom codes.                    //
//                                                                            //
// In in a later version this should be redesigned to give the user           //
// more control over they error system.                                       //
////////////////////////////////////////////////////////////////////////////////
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
