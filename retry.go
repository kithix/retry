package retry

import (
	"math/rand"
	"time"
)

// This takes
func Do(
	doer func() error,
	retryStrategy func(error) bool,
) error {
	for {
		err := doer()
		if err == nil {
			return nil
		}
		if !retryStrategy(err) {
			return err
		}
	}
}

// Always can be used as a strategy to always try again when an error is received.
func Always(_ error) bool {
	return true
}

// Limit wraps a retry strategy to have a limit on retries
func WithLimit(
	retryStrategy func(error) bool,
	limit int,
) func(error) bool {
	i := 0
	return func(err error) bool {
		i++
		if !retryStrategy(err) {
			return false
		}

		if i < limit {
			return true
		}
		return false
	}
}

func Limit(limit int) func(error) bool {
	return WithLimit(Always, limit)
}

// WithWait wraps a retry strategy to pause between retries.
func WithWait(
	retryStrategy func(error) bool,
	timeBetween time.Duration,
) func(error) bool {
	return func(err error) bool {
		if !retryStrategy(err) {
			return false
		}

		time.Sleep(timeBetween)
		return true
	}
}

// WithExponentialBackoff wraps a retry strategy to wait between retries exponentially longer with each one.
// It starts at the minimum and increases at a factor of 'n^2'. It also allows a range of jitter to be added to each one.
func WithExponentialBackoff(
	retryStrategy func(error) bool,
	minimumWait time.Duration,
	maximumWait time.Duration,
	jitter time.Duration,
) func(error) bool {
	retries := 1 * time.Second
	currentWaitTime := minimumWait

	return func(err error) bool {
		if !retryStrategy(err) {
			return false
		}
		retries++
		randomWait := time.Duration(rand.Int63n(int64(jitter)))
		time.Sleep(currentWaitTime + randomWait)

		currentWaitTime = minimumWait * (retries ^ 2)

		if currentWaitTime > maximumWait {
			currentWaitTime = maximumWait
		}

		return true
	}
}
