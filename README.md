# zerr
Error forwarding with stack trace where it is used.

![image info](./vs.png)

Example of usage:
```
package main

import (
	"database/sql"
	"fmt"

	"github.com/zmajew/zerr"
	"log"
)

// Function where the error happened.
func A() error {
	// Database returned this error on query ...
	return zerr.Forward(sql.ErrNoRows)
}

// Some middle function.
func B() error {
	return zerr.Forward(A())
}

func main() {
	err := B()

	// Send an error to the frontend:
	fmt.Println(zerr.WithoutStackTrace(err))

	// Log an error with stack trace:
	log.Fatal(err.Error())
}
```


