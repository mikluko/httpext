package useragent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkString_Any(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		benchmark(pb, Any)
	})
}

func BenchmarkString_Mobile(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		benchmark(pb, Mobile)
	})
}

func benchmark(pb *testing.PB, ua UserAgent) {
	for pb.Next() {
		_ = String(ua)
	}
}

func TestString(t *testing.T) {
	for _, v := range values {
		s := String(v)
		if s == "" {
			t.Fatalf("user agent string is empty: %q", v)
		}
	}
}

func TestString_Rotates(t *testing.T) {
	var (
		prev string
		cur  string
	)
	for i := 0; i < 100; i++ {
		cur = String(Any)
		assert.NotEqual(t, prev, cur)
		prev = cur
	}
}

func TestFromString(t *testing.T) {
	t.Run("camel", func(t *testing.T) {
		for _, v := range values {
			str, ok := stringsCamel[v]
			if !ok {
				t.Errorf("no kebab string for user agent: %v", v)
			}
			ua := FromString(str)
			assert.Equal(t, v, ua)
		}
	})
	t.Run("kebab", func(t *testing.T) {
		for _, v := range values {
			str, ok := stringsKebab[v]
			if !ok {
				t.Errorf("no kebab string for user agent: %v", v)
			}
			ua := FromString(str)
			assert.Equal(t, v, ua)
		}
	})
}
