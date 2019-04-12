package errors

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"fmt"
)

func Test_New(t *testing.T) {
	Convey("new error ", t, func() {

		err := New("test")
		Convey("error has stack trace", func() {
			wt, _ := err.(*withTrace)
			So(wt.stack, ShouldNotBeNil)
		})
	})
}

func Test_WrapMsg(t *testing.T) {
	Convey("new error ", t, func() {

		err := WrapMsg(fmt.Errorf("cause"), "test")
		Convey("error has stack trace", func() {
			wt, _ := err.(*withTrace)
			So(wt.stack, ShouldNotBeNil)
			So(wt.cause, ShouldNotBeNil)
		})
	})
}
