package errorx

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("whoops")
	fmt.Println(err)

	// Output: whoops
}

func TestNew_printf(t *testing.T) {
	err := New("whoops")
	fmt.Printf("%+v", err)

	// Output:
	// whoops
	// command-line-arguments.TestNew_printf
	// 		D:/util/core/errorx/errorx_test.go:16
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
}

func TestOrigNewWithMessage(t *testing.T) {
	cause := errors.New("whoops")
	err := WithMessage(cause, "oh noes")
	fmt.Println(err)
	fmt.Println(err.Error())
	fmt.Println(Cause(err).Error())

	// Output:
	// oh noes: whoops
	// oh noes: whoops
	// whoops
}

func TestWithMessage(t *testing.T) {
	cause := New("whoops")
	err := WithMessage(cause, "oh noes")
	fmt.Println(err)
	fmt.Println(err.Error())
	fmt.Println(Cause(err).Error())

	// Output:
	// oh noes: whoops
	// oh noes: whoops
	// whoops

}

func TestWithStack(t *testing.T) {
	cause := New("whoops")
	err := WithStack(cause)
	fmt.Println(err)

	// Output: whoops
}

func TestWithStack_printf(t *testing.T) {
	cause := New("whoops")
	err := WithStack(cause)
	fmt.Printf("%+v", err)

	// Output:
	// whoops
	// command-line-arguments.TestWithStack_printf
	// 		D:/util/core/errorx/errorx_test.go:54
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
	// command-line-arguments.TestWithStack_printf
	// 		D:/util/core/errorx/errorx_test.go:55
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
}

func TestWrap(t *testing.T) {
	cause := New("whoops")
	err := Wrap(cause, "oh noes")
	fmt.Println(err)

	// Output: oh noes: whoops
}

func fn() error {
	e1 := New("error")
	e2 := Wrap(e1, "inner")
	e3 := Wrap(e2, "middle")
	return Wrap(e3, "outer")
}

func TestCause(t *testing.T) {
	err := fn()
	fmt.Println(err)
	fmt.Println(Cause(err))

	// Output: outer: middle: inner: error
	// error
}

func TestWrap_extended(t *testing.T) {
	err := fn()
	fmt.Printf("%+v\n", err)

	// Output:
	// error
	// command-line-arguments.fn
	// 		D:/util/core/errorx/errorx_test.go:99
	// command-line-arguments.TestWrap_extended
	// 		D:/util/core/errorx/errorx_test.go:115
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
	// inner
	// command-line-arguments.fn
	// 		D:/util/core/errorx/errorx_test.go:100
	// command-line-arguments.TestWrap_extended
	// 		D:/util/core/errorx/errorx_test.go:115
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
	// middle
	// command-line-arguments.fn
	// 		D:/util/core/errorx/errorx_test.go:101
	// command-line-arguments.TestWrap_extended
	// 		D:/util/core/errorx/errorx_test.go:115
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
	// outer
	// command-line-arguments.fn
	// 		D:/util/core/errorx/errorx_test.go:102
	// command-line-arguments.TestWrap_extended
	// 		D:/util/core/errorx/errorx_test.go:115
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
}

func TestWrapf(t *testing.T) {
	cause := New("whoops")
	err := Wrapf(cause, "oh noes #%d", 2)
	fmt.Println(err)

	// Output: oh noes #2: whoops
}

func TestErrorf_extended(t *testing.T) {
	err := Errorf("whoops: %s", "foo")
	fmt.Printf("%+v", err)

	// Output:
	// whoops: foo
	// command-line-arguments.TestErrorf_extended
	// 		D:/util/core/errorx/errorx_test.go:153
	// testing.tRunner
	// 		D:/Go/src/testing/testing.go:1259
	// runtime.goexit
	// 		D:/Go/src/runtime/asm_amd64.s:1581
}

func Test_stackTrace(t *testing.T) {
	type stackTracer interface {
		StackTrace() StackTrace
	}

	err, ok := Cause(fn()).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()
	fmt.Printf("%+v", st[0:2]) // top two frames

	// Output:
	// command-line-arguments.fn
	//     D:/util/core/errorx/errorx_test.go:99
	// command-line-arguments.Test_stackTrace
	// D:/util/core/errorx/errorx_test.go:179
}

func TestCause_printf(t *testing.T) {
	err := Wrap(func() error {
		return func() error {
			return New("hello world")
		}()
	}(), "failed")

	fmt.Printf("%v", err)

	// Output: failed: hello world
}
