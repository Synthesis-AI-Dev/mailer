package generic

import (
	"time"

	"github.com/Synthesis-AI-Dev/mailer/lib/random"
)

// StopErr provides a special return value that the users of Retry can use to
// skip retrying when doing so would be futile.
type StopErr struct {
	error
}

// NewStopErr is a StopErr constructor
func NewStopErr(e error) StopErr {
	return StopErr{e}
}

// Retry provided f any number of times with jittery exponential backoff.
//
// Courtesy Nick via https://upgear.io/blog/simple-golang-retry-function/
func Retry(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(StopErr); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(random.Jitter(sleep))
			return Retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}
