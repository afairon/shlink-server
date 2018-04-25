//go:generate ffjson $GOFILE
package handlers

import (
	"fmt"
)

//go:generate $GOPATH/bin/ffjson $GOFILE
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func httpError(code int, s string, args ...interface{}) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: fmt.Sprintf(s, args...),
	}
}
