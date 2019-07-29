package errors

import "github.com/hashicorp/go-multierror"

// wrap go-mutlierror package
type Error struct {
  multierr *multierror.Error
}

func (e *Error) Error() string {
  return e.multierr.Error()
}


func (e *Error) ErrorOrNil() error {
  if e == nil || e.multierr == nil {
    return nil
  }
  
  return e.multierr.ErrorOrNil()
}

func Append(err error, errs error) *Error {
  if errs == nil {
    return nil
  }

  return &Error {multierror.Append(err, errs)}
}
