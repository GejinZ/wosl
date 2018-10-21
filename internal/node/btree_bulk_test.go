package node

import (
	"fmt"
	"testing"

	"github.com/zeebo/wosl/internal/assert"
)

func TestBtreeBulk(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		var bu btreeBulk
		var buf []byte

		for i := 0; i < 1000; i++ {
			ent, _ := appendEntry(&buf, fmt.Sprint(i), 0)
			bu.append(ent)
		}

		bt := bu.done()

		i := 0
		bt.Iter(func(ent *entry) bool {
			key := string(ent.readKey(buf))
			assert.Equal(t, key, fmt.Sprint(i))
			i++
			return true
		})
	})

	t.Run("One", func(t *testing.T) {
		var bu btreeBulk
		var buf []byte

		ent, _ := appendEntry(&buf, "0", 0)
		bu.append(ent)
		bt := bu.done()

		bt.Iter(func(ent *entry) bool {
			assert.Equal(t, string(ent.readKey(buf)), "0")
			return true
		})
	})

	t.Run("Zero", func(t *testing.T) {
		var bu btreeBulk
		bt := bu.done()

		bt.Iter(func(ent *entry) bool {
			t.Fatal("expected no entries")
			return true
		})
	})
}