package sleepy

import "net/http"

type Filter func(http.ResponseWriter, *http.Request) error
