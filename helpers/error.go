package helpers

import (
	"net/http"
)

var (
	FileNotFound         = "file not found"
	UnsupportedMediaType = "this media type is not supported"
	BadRequest           = "bad request"
)

func GetStatusCodeFromErr(err error) int {
	switch err.Error() {
	case FileNotFound:
		return http.StatusNotFound
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
