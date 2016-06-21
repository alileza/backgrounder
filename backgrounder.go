package backgrounder

import (
	"errors"
	"time"
)

var (
	ErrTimeout = errors.New("Process timeout.")
)

type Backgrounder struct {
	err    chan error
	Errors error
	Count  int
}

func New() *Backgrounder {
	bg := &Backgrounder{}
	bg.err = make(chan error)
	return bg
}

func (bg *Backgrounder) Run(f func() error) {
	bg.Count++
	go func() {
		bg.err <- f()
	}()
}

func (bg *Backgrounder) CatchErr(timeouts ...time.Duration) []error {
	var errs []error
	timeout := time.Minute
	iteration := 1

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	for {
		select {
		case err := <-bg.err:
			if err != nil {
				errs = append(errs, err)
			}
			bg.Count--
		case <-time.After(timeout):
			errs = append(errs, ErrTimeout)
			break
		}
		if iteration >= bg.Count {
			break
		}
		iteration++
	}
	return errs
}
