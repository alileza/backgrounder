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

	err := bg.CatchErrs()
	if len(err) != 0 {
		t.Error("It should be success but it's fail ", err)
	}
	if name != "joni" {
		t.Error("failed to binding global name on background process.")
	}
}

func TestRunTimeout(t *testing.T) {
	bg := New()

	bg.Run(func() error {
		time.Sleep(time.Second * 5000)
		return nil
	})

	err := bg.CatchErrs(time.Millisecond)
	if len(err) == 0 || err[0] != ErrTimeout {
		t.Error("failed to timeout.")
	}
}

func TestCatchError(t *testing.T) {
	bg := New()

	bg.Run(func() error {
		return errors.New("Hello, World")
	})

	err := bg.CatchErrs()
	if err[0].Error() != "Hello, World" {
		t.Error("failed pass error.")
	}
}

func TestCatchAllError(t *testing.T) {
	bg := New()

	bg.Run(func() error {
		return errors.New("-")
	})
	bg.Run(func() error {
		return errors.New("-")
	})
	bg.RunProfile(func() error {
		return errors.New("-")
	}, "3")
	bg.Run(func() error {
		return errors.New("-")
	})

	errs := bg.CatchErrs()
	if bg.Count() != 0 || len(errs) != 4 {
		t.Errorf("Missing %d process.", bg.Count())
	}
}

func TestProfile(t *testing.T) {
	bg := New()

	bg.RunProfile(func() error {
		time.Sleep(time.Millisecond * 2)
		return errors.New("-")
	}, "3")

	bg.CatchErrs()
	if bg.GetProfile("3") < (time.Millisecond * 2) {
		t.Errorf("Expecting profile more then 2 Millisecond, got %v", bg.GetProfile("3"))
	}
}

func TestProfileCount(t *testing.T) {
	bg := New()

	bg.RunProfile(func() error {
		time.Sleep(time.Millisecond * 2)
		return errors.New("-")
	}, "3")
	bg.RunProfile(func() error {
		time.Sleep(time.Millisecond * 2)
		return errors.New("-")
	}, "2")
	bg.RunProfile(func() error {
		time.Sleep(time.Millisecond * 2)
		return errors.New("-")
	}, "2")

	bg.Run(func() error {
		time.Sleep(time.Millisecond * 2)
		return errors.New("-")
	})

	bg.CatchErrs()
	if len(bg.GetProfiles()) != 2 {
		t.Errorf("Expecting profile 2 processes, got %v", len(bg.GetProfiles()))
	}
}

func TestCatchErrorWithNoRun(t *testing.T) {
	bg := New()
	st := time.Now()
	bg.CatchErrs(time.Millisecond)
	if time.Since(st) >= time.Millisecond {
		t.Error("not returning empty process.")
	}
}
