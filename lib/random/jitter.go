package random

import (
	"math/rand"
	"time"
)

// Jitter adds a little random time to the provided duration.
func Jitter(d time.Duration) time.Duration {
	jitter := time.Duration(rand.Int63n(int64(d)))
	return d + jitter/2
}
