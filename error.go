package zuul

import "github.com/tietang/go-utils/errs"

const (
    errNotFoundRouteError = 30001
    errNotFoundShardResource = 30002
    errKeyGreaterThanZero = 30101
    errFallback = 40001
)

func newNotFoundRouteError(code int, message string) *errs.Error {
    return &errs.Error{Message: message, Code: code}
}

func newFallbackError(code int, message string) *errs.Error {
    return &errs.Error{Message: message, Code: code}
}

func _error(innerCode int, innerMessage string, message string) *errs.Error {
    return &errs.Error{Message: innerMessage + ",cause by: " + message, Code: innerCode}
}
