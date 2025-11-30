// модуль errors содержит кастомные ошибки.
package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoContent          = errors.New("no content")
)
