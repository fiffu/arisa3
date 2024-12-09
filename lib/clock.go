package lib

import (
	"testing"
	"time"
)

type frozenClock struct {
	timestamp time.Time
}

func (fc *frozenClock) Now() time.Time {
	return fc.timestamp
}

func (fc *frozenClock) Add(duration time.Duration) {
	fc.timestamp = fc.timestamp.Add(duration)
}

func (fc *frozenClock) Set(timestamp time.Time) {
	fc.timestamp = timestamp
}

func FrozenClock(t *testing.T, timestamp time.Time) *frozenClock {
	return &frozenClock{timestamp}
}

func FrozenNow(t *testing.T) *frozenClock {
	return FrozenClock(t, time.Now())
}
