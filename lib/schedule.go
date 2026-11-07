package lib

import (
	"context"
	"fmt"
	"time"
)

// lib.FrozenClock satisfies Clock in tests.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

// SlotOf uses Unix seconds, so slots line up with the Unix epoch: ts - ts % size.
func SlotOf(t time.Time, size time.Duration) time.Time {
	sec := int64(size / time.Second)
	ts := t.Unix()
	return time.Unix(ts-ts%sec, 0).UTC()
}

func NextSlot(t time.Time, size time.Duration) time.Time {
	return SlotOf(t, size).Add(size)
}

func IsSlot(t time.Time, size time.Duration) bool {
	return t.Unix()%int64(size/time.Second) == 0
}

// Job does the work for one slot. It must be idempotent:
// the runner can call it more than once for the same slot.
type Job func(ctx context.Context, slot time.Time) error

type Runner struct {
	Size  time.Duration
	Job   Job
	Clock Clock
	// Sleep waits for d or until ctx ends. Tests replace it.
	Sleep   func(ctx context.Context, d time.Duration) error
	OnError func(ctx context.Context, slot time.Time, err error)
}

func NewRunner(size time.Duration, job Job) *Runner {
	return &Runner{
		Size:    size,
		Job:     job,
		Clock:   realClock{},
		Sleep:   sleepCtx,
		OnError: func(context.Context, time.Time, error) {},
	}
}

// Run calls Job for the current slot at once, then for each later slot
// when it starts. It returns when ctx ends.
func (r *Runner) Run(ctx context.Context) error {
	slot := SlotOf(r.Clock.Now(), r.Size)
	for {
		r.runSafe(ctx, slot)

		// Compute the next slot before sleeping, so an early timer
		// still runs the right slot.
		slot = slot.Add(r.Size)
		if err := r.Sleep(ctx, slot.Sub(r.Clock.Now())); err != nil {
			return err
		}
	}
}

// runSafe turns a panic into an error. A failed slot does not stop later slots.
func (r *Runner) runSafe(ctx context.Context, slot time.Time) {
	defer func() {
		if p := recover(); p != nil {
			r.OnError(ctx, slot, fmt.Errorf("panic: %v", p))
		}
	}()
	if err := r.Job(ctx, slot); err != nil {
		r.OnError(ctx, slot, err)
	}
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
