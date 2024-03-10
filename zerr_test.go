package zerr_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/zmajew/zerr"
)

func TestForward(t *testing.T) {
	err := sql.ErrNoRows
	ze := zerr.Forward(err)
	ze2 := zerr.Forward(ze)

	if !errors.Is(ze2, err) {
		t.FailNow()
	}

	if !json.Valid([]byte(ze2.Error())) {
		t.FailNow()
	}
}

func TestWithoutStack(t *testing.T) {
	err := sql.ErrNoRows
	ze := fmt.Errorf("wrapped error %w", err)
	ze2 := zerr.Forward(ze)

	ws := zerr.WithoutStack(ze2)

	if ws.Error() != ze.Error() {
		t.FailNow()
	}

	// Check to not unwrap further, since this is not its job.
	wsw := zerr.WithoutStack(ze)
	if wsw.Error() != ze.Error() {
		t.FailNow()
	}
}

func ExampleForward() {
	err := sql.ErrNoRows
	ze := zerr.Forward(err)
	ze2 := zerr.Forward(ze)

	fmt.Fprintln(os.Stdout, zerr.WithoutStack(ze2))

	fmt.Fprintln(os.Stdout, ze2.Error())
}
