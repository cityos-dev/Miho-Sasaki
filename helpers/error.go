package helpers

import (
	"net/http"
)

var (
	ContentTypeIsWrong   = "content type is wrong"
	ParamIsInvalid       = "params is invalid"
	FileNotFound         = "file not found"
	UnsupportedMediaType = "this media type is not supported"
)

func GetStatusCodeFromErr(err error) int {
	switch err.Error() {
	case FileNotFound:
		return http.StatusNotFound
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}
