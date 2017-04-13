package errorcollector

import (
	"fmt"
	"strings"
)

// ErrorCollector is used to collect multiple error
// values while still satisfying the error interface
type ErrorCollector []error

type unwrapper interface {
	unwrap() []error
}

// New returns a new ErrorCollector
func New() ErrorCollector {
	var ec ErrorCollector
	return ec
}

func (ec ErrorCollector) unwrap() []error {
	collection := []error{}
	for _, err := range ec {
		collection = append(collection, err)
	}
	return collection
}

// Collect adds a single error or another collector to
// the collection of errors
func (ec *ErrorCollector) Collect(err error) {
	if err == nil {
		return
	}
	if castCollector, ok := err.(unwrapper); ok {
		*ec = append(*ec, castCollector.unwrap()...)
	} else {
		*ec = append(*ec, err)
	}
}

func (ec ErrorCollector) Error() string {
	collection := []string{}
	for _, err := range ec {
		collection = append(collection, err.Error())
	}
	return fmt.Sprintf("collected errors: %s", strings.Join(collection, ", "))
}
