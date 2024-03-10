package zerr

import (
	"fmt"
	"runtime"
)

type zError struct {
	stackTrace []string
	err        error
}

// Error defines zError as error type.
func (e zError) Error() string {
	// add last location to stack
	fileName, line, funcName := location(3)
	e.stackTrace = append(e.stackTrace, fmt.Sprintf(`{"location": "%s:%d", "function": "%s"}`, fileName, line, funcName))

	// create json error string
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

// Forward picks stack trace of the error it receives.
func Forward(err error) error {
	if err != nil {
		fileName, line, funcName := location(2)
		errorLocation := fmt.Sprintf(`{"location": "%s:%d", "function": "%s"}`, fileName, line, funcName)

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

func location(level int) (string, int, string) {
	pc := make([]uintptr, 10)
	runtime.Callers(level, pc)
	f := runtime.FuncForPC(pc[1] - 1)
	_, fn, line, _ := runtime.Caller(level)
	return fn, line, f.Name()
}

// Unwrap adds an unwrap functionality to zError type.
func (c zError) Unwrap() error {
	return c.err
}

// WithoutStack looks for first form backwards zError type an returns its encapsulated error.
// If zError was not found it will return an error from argument.
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
