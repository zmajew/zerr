package zerr

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

type EstudyError struct {
	ErrorLocation string
	Err           error
}

func (e *EstudyError) Error() string {
	return fmt.Sprintf("%s\n%s", e.Err.Error(), e.ErrorLocation)
}

func Forward(err error) *EstudyError {
	osPath, _ := os.Getwd()
	estudy := &EstudyError{}

	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1] - 1)

	estudy.Err = err
	_, fn, line, _ := runtime.Caller(1)
	fn = strings.TrimPrefix(fn, osPath)
	errorLocation := fmt.Sprintf("%s %d, %s", fn, line, f.Name())
	estudy.ErrorLocation = errorLocation
	return estudy
}

func ForwardWithMessage(err error, text string) *EstudyError {
	osPath, _ := os.Getwd()
	estudy := &EstudyError{}

	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1] - 1)

	estudy.Err = err
	_, fn, line, _ := runtime.Caller(1)
	fn = strings.TrimPrefix(fn, osPath)
	errorLocation := fmt.Sprintf("%s: %s %d, %s", text, fn, line, f.Name())
	estudy.ErrorLocation = errorLocation

	return estudy
}

func (c *EstudyError) unwrap() error {
	return c.Err
}

func GetFirstError(err error) error {
	for {
		b, ok := err.(*EstudyError)
		if !ok {
			return err
		}
		err = b.unwrap()
	}
}

func Log(err error) {
	osPath, _ := os.Getwd()

	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1] - 1)

	_, fn, line, _ := runtime.Caller(1)
	fn = strings.TrimPrefix(fn, osPath)
	location := fmt.Sprintf("%s %d, %s", fn, line, f.Name())

	t := time.Now()

	// Yellow error color warn
	yellError := color.New(color.FgYellow).Add(color.Bold)

	// Print the message
	yellError.Printf("Error:\n")
	fmt.Printf("Time: %s\n", t.String())
	fmt.Println(err.Error())
	fmt.Println(location)
}
