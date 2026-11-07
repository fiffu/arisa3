package lib

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const fourHours = 4 * time.Hour

func Test_SlotOf(t *testing.T) {
	in := time.Unix(1762510000, 0) // 2025-11-07 10:06:40 UTC
	assert.Equal(t, int64(1762502400), SlotOf(in, fourHours).Unix())
	assert.Equal(t, int64(1762502400+4*3600), NextSlot(in, fourHours).Unix())

	// A slot start maps to itself.
	assert.Equal(t, int64(1762502400), SlotOf(time.Unix(1762502400, 0), fourHours).Unix())
}

func Test_IsSlot(t *testing.T) {
	assert.True(t, IsSlot(time.Unix(1762502400, 0), fourHours))
	assert.False(t, IsSlot(time.Unix(1762502401, 0), fourHours))
}

func Test_Runner_runsCurrentSlotThenEachNextSlot(t *testing.T) {
	start := time.Unix(1762510000, 0)
	clock := FrozenClock(t, start)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var got []int64
	r := NewRunner(fourHours, func(_ context.Context, slot time.Time) error {
		got = append(got, slot.Unix())
		if len(got) == 3 {
			cancel()
		}
		return nil
	})
	r.Clock = clock
	r.Sleep = func(ctx context.Context, d time.Duration) error {
		clock.Add(d)
		return ctx.Err()
	}

	err := r.Run(ctx)

	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, []int64{1762502400, 1762502400 + 14400, 1762502400 + 28800}, got)
}

func Test_Runner_errorsAndPanicsDoNotStopLaterSlots(t *testing.T) {
	clock := FrozenClock(t, time.Unix(1762502400, 0))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	var errs []error
	r := NewRunner(fourHours, func(context.Context, time.Time) error {
		calls++
		switch calls {
		case 1:
			return errors.New("boom")
		case 2:
			panic("kaboom")
		default:
			cancel()
			return nil
		}
	})
	r.Clock = clock
	r.Sleep = func(ctx context.Context, d time.Duration) error {
		clock.Add(d)
		return ctx.Err()
	}
	r.OnError = func(_ context.Context, _ time.Time, err error) { errs = append(errs, err) }

	_ = r.Run(ctx)

	assert.Equal(t, 3, calls)
	assert.Len(t, errs, 2)
	assert.EqualError(t, errs[1], "panic: kaboom")
}
