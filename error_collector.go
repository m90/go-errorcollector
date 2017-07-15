// Package errorcollector eases handling of
// multiple errors in the same context
package errorcollector

import (
	"fmt"
	"strings"
)

// ErrorCollector is used to collect multiple error
// values while satisfying the common error interface
type ErrorCollector []error

// New returns a new nil ErrorCollector
func New() ErrorCollector {
	var ec ErrorCollector
	return ec
}

// Collect adds a single error or another collector to
// the collection of errors, passing nil is a noop
func (ec *ErrorCollector) Collect(err error) {
	if err == nil {
		return
	}
	if castCollector, ok := err.(ErrorCollector); ok {
		*ec = append(*ec, castCollector...)
	} else {
		*ec = append(*ec, err)
	}
}

// Error returns a string describing the collected errors
func (ec ErrorCollector) Error() string {
	if len(ec) == 1 {
		return ec[0].Error()
	}
	collection := []string{}
	for _, err := range ec {
		collection = append(collection, err.Error())
	}
	return fmt.Sprintf("collected errors: %s", strings.Join(collection, ", "))
}
