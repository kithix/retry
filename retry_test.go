package retry

import (
	"errors"
	"testing"
)

var testError = errors.New("test error")
var testDontRetryError = errors.New("dont retry this error")

func errorForNCalls(n int) func() error {
	i := 0
	return func() error {
		if n < 0 {
			return testError
		}
		i++
		if i > n {
			return nil
		}
		return testError
	}
}

var errorForever = errorForNCalls(-1)

func TestRetryLimit(t *testing.T) {
	runCount := 0
	err := Do(
		func() error {
			runCount++
			return errorForever()
		},
		Limit(5),
	)
	if err != testError {
		t.Error("Unexpected error", err)
	}
	if runCount != 5 {
		t.Error("Unexpected run count, got", runCount)
	}
}

func TestRetryStrategy(t *testing.T) {
	runCount := 0
	err := Do(
		func() error {
			runCount++
			if runCount == 5 {
				return nil
			}
			return errorForever()
		},
		Always,
	)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if runCount != 5 {
		t.Error("Unexpected run count, got", runCount)
	}
}

func TestErrorEvaluation(t *testing.T) {
	runCount := 0
	err := Do(
		func() error {
			runCount++
			if runCount < 3 {
				return testError
			}
			return testDontRetryError
		},
		func(err error) bool {
			if err == testDontRetryError {
				return false
			}
			return true
		},
	)
	if err != testDontRetryError {
		t.Error("Unexpected error", err)
	}
	if runCount != 3 {
		t.Error("Unexpected run count, got", runCount)
	}
}
