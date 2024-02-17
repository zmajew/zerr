package zerr

import (
	"fmt"
	"runtime"
)

type zError struct {
	stackTrace []string
	err        error
}

func (e zError) Error() string {
	a := "{"
	b := fmt.Sprintf(`"error": "%s", `, e.err.Error())
	c := `"stack_trace": [`
	for i, v := range e.stackTrace {
		if i == 0 {
			c = c + v
			continue
		}
		c = c + "," + v
	}
	c = c + "]"

	return a + b + c + "}"
}

func Forward(err error) error {
	if err != nil {
		pc := make([]uintptr, 10)
		runtime.Callers(1, pc)
		f := runtime.FuncForPC(pc[1] - 1)
		_, fn, line, _ := runtime.Caller(1)
		errorLocation := fmt.Sprintf(`{"location": "%s:%d", "function": "%s"}`, fn, line, f.Name())

		zerr, ok := err.(zError)
		if ok {
			zerr.stackTrace = append(zerr.stackTrace, errorLocation)
			return zerr
		}
		return zError{
			err:        err,
			stackTrace: append(zerr.stackTrace, errorLocation),
		}
	}

	return nil
}

// Unwrap adds an unwrap standard error functionality.
func (c zError) Unwrap() error {
	return c.err
}

// WithoutStack looks for first form backwards zError type an returns its encapsulated error.
// If does not found zError with Unwrap function will return an error from argument.
func WithoutStack(err error) error {
	if ze, ok := err.(zError); ok {
		return ze.err
	}

	if x, ok := err.(interface{ Unwrap() error }); ok {
		if _, ok := err.(zError); !ok {
			return err
		}
		return WithoutStack(x.Unwrap())
	}

	return err
}
