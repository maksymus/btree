package errors

import (
  "bytes"
  "fmt"
  "runtime"
  "strings"
)

// error code with trace
type withTrace struct {
  msg   string
  stack string
  cause error
}

func (err *withTrace) Error() string {
  return err.msg
}

// new error with trace
func New(msg string) error {
  return &withTrace{
    msg:   msg,
    stack: trace(),
  }
}

// wrap existing error adding trace
func Wrap(err error) error {
  return &withTrace{
    msg:   err.Error(),
    stack: trace(),
    cause: err,
  }
}

// wrap existing error adding message and trace
func WrapMsg(err error, msg string) error {
  return &withTrace{
    msg:   msg,
    stack: trace(),
    cause: err,
  }
}

func Stack(err error) string {
  if wt, ok := err.(*withTrace); ok {
    return wt.stack
  }

  return ""
}

func Cause(err error) error {
  if wt, ok := err.(*withTrace); ok {
    return wt.cause
  }

  return nil
}

// Get stack frame pointers
func frames(depth int) []uintptr {
  pcs := make([]uintptr, depth)
  n := runtime.Callers(5, pcs[:])
  return pcs[0:n]
}

// Formats stack trace and returns resulting string
func trace() string {
  var buf bytes.Buffer
  traceWalk(16, func(fn function) {
    fmt.Fprintf(&buf, "%s\n\t%s:%d\n", fn.name, fn.file, fn.line)
  })
  return buf.String()
}

func traceWalk(depth int, callback func(fn function)) {
  frames := frames(depth)
  for _, pointer := range frames {
    fr := frame(pointer)

    f, err := fr.stat()
    if err == nil {
      callback(f)
    }
  }
}

// Method frame type
type frame uintptr
type function struct {
  name string
  file string
  line int
}

// Frame pointer
func (f frame) pointer() uintptr {
  return uintptr(f) - 1
}

// Frame stat
func (f frame) stat() (function, error) {
  fn := runtime.FuncForPC(f.pointer())
  if fn == nil {
    return function{}, fmt.Errorf("FuncForPC returned nil")
  }

  // last index of '/' to truncate function name
  lastIndex := strings.LastIndex(fn.Name(), "/")
  file, line := fn.FileLine(f.pointer())
  return function{fn.Name()[lastIndex+1:], file, line}, nil
}
