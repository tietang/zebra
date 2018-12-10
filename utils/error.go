package utils

import "github.com/tietang/go-utils/errs"

const (
    ErrNotFoundRouteError    = 30001
    ErrNotFoundShardResource = 30002
    ErrKeyGreaterThanZero    = 30101
    ErrFallback              = 40001
)

func NewNotFoundRouteError(code int, message string) *errs.Error {
    return &errs.Error{Message: message, Code: code}
}

func NewFallbackError(code int, message string) *errs.Error {
    return &errs.Error{Message: message, Code: code}
}

func _error(innerCode int, innerMessage string, message string) *errs.Error {
    return &errs.Error{Message: innerMessage + ",cause by: " + message, Code: innerCode}
}
