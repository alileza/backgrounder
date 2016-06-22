package backgrounder

import (
	"errors"
	"time"
)

// Standard error list.
var (
	ErrTimeout = errors.New("Process timeout.")
)

// Handler is func abstraction.
type Handler func() error

type process struct {
	Name        string
	ProcessTime time.Duration
	Error       error
}

// Backgrounder is the main module.
type Backgrounder struct {
	pipe    chan process
	errors  []error
	count   int
	timeout bool

	profile map[string]time.Duration
}

// New returns backgrounder object.
func New() *Backgrounder {
	bg := &Backgrounder{}
	bg.pipe = make(chan process)

	bg.profile = make(map[string]time.Duration)
	return bg
}

// Run accept function, that will automatically executed
// on go routines.
func (bg *Backgrounder) Run(f Handler) {
	bg.count++
	go func() {
		bg.pipe <- process{Error: f()}
	}()
}

// RunProfile accept function, that will automatically executed
// on go routines.
func (bg *Backgrounder) RunProfile(f Handler, name string) {
	bg.count++
	go func() {
		startTime := time.Now()
		bg.pipe <- process{
			Name:        name,
			Error:       f(),
			ProcessTime: time.Since(startTime),
		}
	}()
}

// GetProfile returns profiling result.
func (bg *Backgrounder) GetProfile(key string) time.Duration {
	return bg.profile[key]
}

// GetProfiles returns profiling result.
func (bg *Backgrounder) GetProfiles() map[string]time.Duration {
	return bg.profile
}

// Count returns how many background process
// that currently running.
func (bg *Backgrounder) Count() int {
	return bg.count
}

// Errs returns error list of background process
// that already finished.
func (bg *Backgrounder) Errs() []error {
	return bg.errors
}

// CatchErrs returns error list of background process
// that runs previously.
// Timeouts parameter is optional, default timeout is `1m`
func (bg *Backgrounder) CatchErrs(timeouts ...time.Duration) []error {
	var errs []error
	if bg.count == 0 {
		return errs
	}
	timeout := time.Minute

	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	for {
		select {
		case data := <-bg.pipe:
			if data.Error != nil {
				errs = append(errs, data.Error)
			}
			bg.count--

			if data.ProcessTime != 0 {
				bg.profile[data.Name] = data.ProcessTime
			}
		case <-time.After(timeout):
			errs = append(errs, ErrTimeout)
			bg.timeout = true
			break
		}
		if bg.count == 0 || bg.timeout {
			break
		}
	}

	bg.errors = errs
	return errs
}
