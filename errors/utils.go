package errors

import "github.com/hashicorp/go-multierror"

// wrap go-mutlierror package
type Error struct {
  *multierror.Error
}

func Append(err error, errs ...error) *Error {
  return &Error {multierror.Append(err, errs...)}
}