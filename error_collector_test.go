package errorcollector

import (
	"errors"
	"testing"
)

func TestErrorCollector(t *testing.T) {
	tests := []struct {
		funcs         []func() error
		notNil        bool
		expectedError string
	}{
		{
			[]func() error{},
			false,
			"",
		},
		{
			[]func() error{
				func() error { return nil },
				func() error { return nil },
				func() error { return nil },
			},
			false,
			"",
		},
		{
			[]func() error{
				func() error { return errors.New("beep") },
				func() error { return nil },
				func() error { return nil },
			},
			true,
			"collected errors: beep",
		},
		{
			[]func() error{
				func() error { return errors.New("beep") },
				func() error { return nil },
				func() error { return errors.New("boop") },
			},
			true,
			"collected errors: beep, boop",
		},
		{
			[]func() error{
				func() error {
					collector := New()
					subCollector := New()
					collector.Collect(errors.New("beep"))
					collector.Collect(errors.New("boop"))
					subCollector.Collect(errors.New("biip"))
					subCollector.Collect(nil)
					collector.Collect(subCollector)
					return collector
				},
				func() error { return nil },
				func() error { return errors.New("baap") },
			},
			true,
			"collected errors: beep, boop, biip, baap",
		},
		{
			[]func() error{
				func() error {
					collector := New()
					subCollector := New()
					collector.Collect(subCollector)
					return collector
				},
				func() error { return nil },
				func() error { return nil },
			},
			false,
			"",
		},
	}

	for _, test := range tests {
		errors := New()
		for _, fn := range test.funcs {
			errors.Collect(fn())
		}
		if (errors != nil) != test.notNil {
			t.Errorf(
				"Expected error to be %v, error was %v",
				test.notNil,
				errors)
		}
		if errors != nil && errors.Error() != test.expectedError {
			t.Errorf(
				"Expected error to return %v, got %v",
				test.expectedError,
				errors.Error())
		}

	}
}

func BenchmarkErrorCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		collector := New()
		childCollector := New()
		collector.Collect(errors.New("beep"))
		childCollector.Collect(errors.New("boop"))
		collector.Collect(childCollector)
	}
}
