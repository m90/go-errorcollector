package errorcollector

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrorCollector(t *testing.T) {
	tests := []struct {
		name          string
		funcs         []func() error
		notNil        bool
		expectedError string
	}{
		{
			"noop",
			[]func() error{},
			false,
			"",
		},
		{
			"no errors",
			[]func() error{
				func() error { return nil },
				func() error { return nil },
				func() error { return nil },
			},
			false,
			"",
		},
		{
			"single error",
			[]func() error{
				func() error { return errors.New("beep") },
				func() error { return nil },
				func() error { return nil },
			},
			true,
			"beep",
		},
		{
			"nested collectors",
			[]func() error{
				func() error { return nil },
				func() error {
					collector := New()
					collector.Collect(errors.New("beep"))
					return collector
				},
				func() error { return nil },
			},
			true,
			"beep",
		},
		{
			"multiple errors flat",
			[]func() error{
				func() error { return errors.New("beep") },
				func() error { return nil },
				func() error { return errors.New("boop") },
			},
			true,
			"collected errors: beep, boop",
		},
		{
			"nested multiple levels",
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
			"joint empty collectors",
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
		t.Run(test.name, func(t *testing.T) {
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
		})
	}
}

func BenchmarkErrorCollector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		collector := New()
		childCollector := New()
		collector.Collect(errors.New("beep"))
		childCollector.Collect(errors.New("boop"))
		collector.Collect(childCollector)
		collector.Error()
	}
}

func ExampleErrorCollector() {
	makeLowerCase := func(str string) (string, error) {
		if strings.ToLower(str) != str {
			return strings.ToLower(str), fmt.Errorf("string %v wasn't all lowercase", str)
		}
		return str, nil
	}
	list := []string{"beep", "boOp", "Baap"}
	result := []string{}
	err := New()
	for _, str := range list {
		lowercased, lcErr := makeLowerCase(str)
		err.Collect(lcErr)
		result = append(result, lowercased)
	}
	if err != nil {
		fmt.Printf("got error: %v\n", err)
	}
	fmt.Printf("lowercased strings: %v", strings.Join(result, ", "))
	// Output:
	// got error: collected errors: string boOp wasn't all lowercase, string Baap wasn't all lowercase
	// lowercased strings: beep, boop, baap
}
