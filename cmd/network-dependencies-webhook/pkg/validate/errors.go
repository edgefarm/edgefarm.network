package validate

import "errors"

var ErrNotFound = errors.New("not found")
var ErrMultipleFound = errors.New("multiple found")
var ErrInvalid = errors.New("invalid")
