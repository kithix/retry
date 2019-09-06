package retry

import (
	"errors"
	"testing"
)

var ErrTest = errors.New("test error")
var ErrTestDontRetry = errors.New("dont retry this error")

func errorForNCalls(n int) func() error {
	i := 0
	return func() error {
		if n < 0 {
			return ErrTest
		}
		i++
		if i > n {
			return nil
		}
		return ErrTest
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
	if err != ErrTest {
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
				return ErrTest
			}
			return ErrTestDontRetry
		},
		func(err error) bool {
			if err == ErrTestDontRetry {
				return false
			}
			return true
		},
	)
	if err != ErrTestDontRetry {
		t.Error("Unexpected error", err)
	}
	if runCount != 3 {
		t.Error("Unexpected run count, got", runCount)
	}
}
