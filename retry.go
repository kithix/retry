package retry

import (
	"math/rand"
	"time"
)

// Do takes a function that can fail with a retry strategy for how to handle that failure
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

// Always will continously retry until there is no error.
func Always(_ error) bool {
	return true
}

// WithLimit wraps a retry strategy to have a limit on the amount of retries
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

// Limit is a helper function to always retry until a count is reached.
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
// It starts at the minimum and increases at a factor of 'n^2'.
// It also allows a range of jitter to be added to each one.
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
