package gpio

import "errors"

var (
	ErrInvalidDirection = errors.New("invalid direction")
	ErrAlreadyExported  = errors.New("already exported")
	ErrNotExported      = errors.New("not exported")
	ErrUnknown          = errors.New("unknown error")
	ErrInvalidPin       = errors.New("invalid pin")
	ErrNotImplemented   = errors.New("not implemented")
)
