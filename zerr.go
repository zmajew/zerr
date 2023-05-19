package zerr

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ZError struct {
	ErrorLocation string
	Err           error
}

func (e ZError) Error() string {
	return fmt.Sprintf("%s\n%s", e.Err.Error(), e.ErrorLocation)
}

func Forward(err error) error {
	if err != nil {
		osPath, _ := os.Getwd()
		zErr := ZError{}

		pc := make([]uintptr, 10)
		runtime.Callers(1, pc)
		f := runtime.FuncForPC(pc[1] - 1)

		zErr.Err = err
		_, fn, line, _ := runtime.Caller(1)
		fn = strings.TrimPrefix(fn, osPath)
		errorLocation := fmt.Sprintf("%s %d, %s", fn, line, f.Name())
		zErr.ErrorLocation = errorLocation

		return zErr
	}

	return nil
}

func ForwardWithMessage(err error, text string) error {
	osPath, _ := os.Getwd()
	zErr := ZError{}

	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1] - 1)

	zErr.Err = err
	_, fn, line, _ := runtime.Caller(1)
	fn = strings.TrimPrefix(fn, osPath)
	errorLocation := fmt.Sprintf("%s: %s %d, %s", text, fn, line, f.Name())
	zErr.ErrorLocation = errorLocation

	return zErr
}

func (c *ZError) unwrap() error {
	return c.Err
}

func GetFirstError(err error) error {
	for {
		b, ok := err.(ZError)
		if !ok {
			return err
		}
		err = b.unwrap()
	}
}

func Log(err error, before ...int) {
	step := 1
	if len(before) != 0 {
		step = before[0]
	}

	osPath, _ := os.Getwd()

	pc := make([]uintptr, 10)
	runtime.Callers(1, pc)
	f := runtime.FuncForPC(pc[1] - 1)

	_, fn, line, _ := runtime.Caller(step)
	fn = strings.TrimPrefix(fn, osPath)
	location := fmt.Sprintf("%s %d, %s", fn, line, f.Name())

	t := time.Now()

	// Yellow error color warn
	yellError := color.New(color.FgYellow).Add(color.Bold)

	// Print the message
	yellError.Printf("Error:\n")
	fmt.Printf("Time: %s\n", t.Format("Mon Jan _2 15:04:05 2006"))
	fmt.Println(err.Error())
	fmt.Println(location)
}
