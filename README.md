# zerr
Error forwarding package

Example of usage:
```
package main

import (
	"database/sql"
	"fmt"

	"github.com/zmajew/zerr"
)

// Function where the error happened
func A() error {
	// Database returned this error on query with fictional id...
	err := sql.ErrNoRows

	// Forward the error
	return zerr.Forward(err)
}

// Some middle function
func B() error {
	err := A()

	return zerr.Forward(err)
}

// Add a comment to the passing error
func C() error {
	err := B()

	return zerr.ForwardWithMessage(err, "some error message")
}

func main() {
	err := C()

	if zerr.GetFirstError(err) == sql.ErrNoRows {
		// Send the error to the frontend
		fmt.Println("There is no rows with requested id in the database")
	}

	zerr.Log(err)
}
```
