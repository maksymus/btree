package errors

// error code with trace
type withTrace struct {

}

func (withTrace) Error() string {
  panic("implement me")
}

// new error with trace
func New(msg string) error {
  return nil
}

// wrap existing error adding trace
func Wrap(err error) error {
  return nil
}

// wrap existing error adding message and trace
func WrapMsg(err error, msg string) error {
  return nil
}

