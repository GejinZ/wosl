package node

import (
	"fmt"
	"testing"

	"github.com/zeebo/wosl/internal/assert"
)

func TestNodeBulk(t *testing.T) {
	t.Run("Append", func(t *testing.T) {
		var bu Bulk

		for i := 0; i < 1000; i++ {
			key := []byte(fmt.Sprintf("%04d", i))
			assert.That(t, bu.Append(key, nil, false))
		}
		n := bu.Done(0, 0)

		last := ""
		n.iter(func(ent *Entry, buf []byte) bool {
			key := string(ent.readKey(buf))
			assert.That(t, key > last)
			last = key
			return true
		})
	})
}

func BenchmarkNodeBulk(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		run := func(b *testing.B, v []byte) {
			keys := make([][]byte, b.N)
			for i := range keys {
				keys[i] = []byte(fmt.Sprintf("%08d", i))
			}

			var bu Bulk
			b.SetBytes(8 + int64(len(v)))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				if !bu.Fits(keys[i], v, bufferSize) {
					bu.Reset()
				}
				bu.Append(keys[i], v, true)
			}
		}

		b.Run("32KB", func(b *testing.B) { run(b, megabuf) })
		b.Run("1KB", func(b *testing.B) { run(b, megabuf[:1<<10]) })
		b.Run("16B", func(b *testing.B) { run(b, megabuf[:1<<4]) })
		b.Run("0B", func(b *testing.B) { run(b, nil) })
	})
}