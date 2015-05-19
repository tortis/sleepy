package sleepy

type Error int

const (
	begin_Internal = iota
	ERR_INTERNAL
	ERR_DATABASE
	ERR_UNKNOWN
	end_Internal

	begin_Authentication
	ERR_INVALID_CREDS
	ERR_EXPIRED_TOKEN
	ERR_INVALID_TOKEN
	end_Authentication

	begin_Request

	end_Request
)

func (e Error) isInternal() bool {
	if e > begin_Internal && e < end_Internal {
		return true
	}
	return false
}
