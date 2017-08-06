package errorcollector_test

import (
	"errors"
	"fmt"
	"testing"

	errorcollector "github.com/m90/go-errorcollector"
)

func TestErrorCollector(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		err := errorcollector.New()
		if err != nil {
			t.Error("Empty collector is not nil")
		}
		err.Collect(nil)
		if err != nil {
			t.Error("Collector is not nil")
		}
		err.Collect(errors.New("this is a test error"))
		if err == nil {
			t.Error("Collector with errors is nil")
		}
		if err.Error() != "this is a test error" {
			t.Errorf("Unexpected error message %v", err.Error())
		}
		err.Collect(errors.New("another one"))
		if err == nil {
			t.Error("Collector with errors is nil")
		}
		if err.Error() != "collected errors: this is a test error, another one" {
			t.Errorf("Unexpected error message %v", err.Error())
		}
	})
	t.Run("error returning func", func(t *testing.T) {
		tester := func(nums ...int) error {
			err := errorcollector.New()
			for _, num := range nums {
				if num%2 != 0 {
					err.Collect(fmt.Errorf("%v is not an even number", num))
				}
			}
			if err != nil {
				return err
			}
			return nil
		}
		if err := tester(2, 4, 6); err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if err := tester(1, 3, 5); err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
