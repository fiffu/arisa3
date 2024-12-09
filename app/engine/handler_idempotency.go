package engine

import (
	"sync"
	"time"

	"github.com/google/btree"
)

// Requests that share an idempotencyKey with another request within this time window
// will be considered duplicates which should be ignored.
const idempotencyWindow = 2 * time.Minute

type idempotencyKey struct {
	key        string
	acquiredAt time.Time
}

func (ik *idempotencyKey) Before(other *idempotencyKey) bool {
	return ik.acquiredAt.Before(other.acquiredAt)
}

type idempotency struct {
	tree  *btree.BTreeG[*idempotencyKey]
	mutex *sync.Mutex
	clock func() time.Time
}

func newIdempotencyChecker() *idempotency {
	var mu sync.Mutex
	var degree int = 2 // idk wtf this means, looks like higher = lower memory usage but slower inserts?
	return &idempotency{
		tree:  btree.NewG(degree, func(a, b *idempotencyKey) bool { return a.Before(b) }),
		mutex: &mu,
		clock: time.Now,
	}
}

// Check attempts to acquire a lock against the given key.
// If there is no existing lock, it creates the lock using the given key and retuns true.
// Otherwise, it returns false.
func (hi *idempotency) Check(key string) (acquired bool) {
	hi.mutex.Lock()
	defer hi.mutex.Unlock()

	now := hi.clock()
	acquired = true

	hi.tree.Ascend(func(item *idempotencyKey) bool {
		age := now.Sub(item.acquiredAt)
		if age > idempotencyWindow {
			hi.tree.Delete(item)
			return true // continue iterating tree
		}

		if item.key == key {
			acquired = false
			return false
		}
		return true
	})

	if acquired {
		hi.tree.ReplaceOrInsert(&idempotencyKey{
			key:        key,
			acquiredAt: now,
		})
	}
	return
}
