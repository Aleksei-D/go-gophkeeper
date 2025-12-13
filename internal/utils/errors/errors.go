// модуль errors содержит кастомные ошибки.
package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoContent          = errors.New("no content")
	ErrCardNumber         = errors.New("card number is invalid")
	ErrExpirationMonth    = errors.New("expiration month is invalid")
	ErrExpirationYear     = errors.New("expiration year is invalid")
	ErrCVV                = errors.New("CVV is invalid")
	ErrCardAlreadyExists  = errors.New("card already exists")
)
