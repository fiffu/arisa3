package colours

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_migrationsDir(t *testing.T) {
	here := migrationsDir
	// assert that this path appears in this order in `migrationsDir`
	expectPath := []string{"app", "cogs", "colours", "dbmigrations"}
	prev := -1
	for _, expect := range expectPath {
		idx := strings.Index(here, expect)
		assert.Less(t, prev, idx)
		prev = idx
	}
}
