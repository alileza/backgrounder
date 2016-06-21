package backgrounder

import (
	"errors"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	var name string
	bg := New()

	bg.Run(func() error {
		time.Sleep(time.Millisecond * 500)
		name = "joni"
		return nil
	})

	err := bg.CatchErr()
	if len(err) != 0 {
		t.Error("It should be success but it's fail ", err)
	}
	if name != "joni" {
		t.Error("failed to binding global name on background process.")
	}
}

func TestRunTimeout(t *testing.T) {
	var name string
	bg := New()

	bg.Run(func() error {
		time.Sleep(time.Second * 5000)
		name = "joni"
		return nil
	})

	err := bg.CatchErr(time.Millisecond)
	if len(err) == 0 || err[0] != ErrTimeout {
		t.Error("failed to timeout.")
	}
}

func TestCatchError(t *testing.T) {
	bg := New()

	bg.Run(func() error {
		return errors.New("Hello, World")
	})

	err := bg.CatchErr()
	if err[0].Error() != "Hello, World" {
		t.Error("failed pass error.")
	}
}
