package helpers

import (
	"net/http"
)

var (
	ContentTypeIsWrong = "content type is wrong"
	ParamIsInvalid     = "params is invalid"
	FileNotFound       = "file not found"
)

func GetStatusCodeFromErr(err error) int {
	switch err.Error() {
	case FileNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
